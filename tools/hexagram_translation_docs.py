#!/usr/bin/env python3
import argparse, json, re
from html import escape, unescape
from pathlib import Path

TEXT_FIELDS=[('name','name'),('title','title'),('above.symbolic','above.symbolic'),('above.alchemical','above.alchemical'),('below.symbolic','below.symbolic'),('below.alchemical','below.alchemical'),('judgment.text','judgment.text'),('judgment.comments','judgment.comments'),('image.text','image.text'),('image.comments','image.comments')]
for i in range(1,7):
    TEXT_FIELDS.append((f'lines.{i}.text',f'lines.{i}.text'))
    TEXT_FIELDS.append((f'lines.{i}.comments',f'lines.{i}.comments'))
SECTION_RE=re.compile(r'^@@HEX\|(\d+)\|([^|]+)\|START@@$')
END_RE=re.compile(r'^@@HEX\|(\d+)\|([^|]+)\|END@@$')

def get_path(obj,path,default=''):
    cur=obj
    for part in path.split('.'):
        if isinstance(cur,dict): cur=cur.get(part,default)
        else: return default
    return cur if cur is not None else default

def load_json(path): return json.loads(Path(path).read_text(encoding='utf-8'))

def export_txt(en_json,out_txt):
    data=load_json(en_json); parts=['PŘEKLADOVÝ SOUBOR I-ŤING\n','Překládej pouze text mezi START/END značkami. Značky samotné neměň.\n']
    for item in data:
        parts.append(f"\n# HEXAGRAM {item['number']} — {item.get('title','')}\n")
        parts.append(f"# {item.get('character','')} {item.get('traditional','')} {item.get('pinyin','')}\n")
        for _,path in TEXT_FIELDS:
            source=get_path(item,path,'')
            parts.append(f"@@HEX|{item['number']}|{path}|START@@\n{source}\n@@HEX|{item['number']}|{path}|END@@\n")
    Path(out_txt).write_text(''.join(parts),encoding='utf-8')

def export_html(en_json,out_html):
    data=load_json(en_json); rows=[]
    for item in data:
        rows.append(f"<h2>Hexagram {item['number']} — {escape(item.get('title',''))}</h2>")
        rows.append(f"<p><strong>{escape(item.get('character',''))}</strong> {escape(item.get('traditional',''))} {escape(item.get('pinyin',''))}</p>")
        for _,path in TEXT_FIELDS:
            source=get_path(item,path,'')
            rows.append(f"<div class='block'><div class='meta'>{item['number']} | {escape(path)}</div><pre>@@HEX|{item['number']}|{path}|START@@\n{escape(source)}\n@@HEX|{item['number']}|{path}|END@@</pre></div>")
    html=f"<!doctype html><html lang='cs'><head><meta charset='utf-8'><title>I-ťing překlad</title><style>body{{font-family:system-ui,sans-serif;max-width:960px;margin:0 auto;padding:24px;line-height:1.5}}pre{{white-space:pre-wrap;background:#f6f6f6;padding:12px;border-radius:8px}}.meta{{font:600 14px system-ui;color:#555;margin-bottom:6px}}.block{{margin:0 0 16px}}</style></head><body><h1>Překladový soubor I-ťing</h1><p>Překládej pouze text mezi značkami START/END. Samotné značky neměň.</p>{''.join(rows)}</body></html>"
    Path(out_html).write_text(html,encoding='utf-8')

def html_to_text(raw:str)->str:
    raw=unescape(raw)
    raw=re.sub(r'<br\s*/?>','\n',raw,flags=re.I)
    raw=re.sub(r'</(p|div|h1|h2|h3|pre|li|tr|td|section|article)>','\n',raw,flags=re.I)
    raw=re.sub(r'<[^>]+>','',raw)
    return raw

def read_translation_text(path:Path)->str:
    suf=path.suffix.lower()
    if suf in {'.html','.htm'}:
        return html_to_text(path.read_text(encoding='utf-8',errors='ignore'))
    if suf=='.txt':
        return path.read_text(encoding='utf-8',errors='ignore')
    if suf=='.docx':
        import zipfile
        from xml.etree import ElementTree as ET
        with zipfile.ZipFile(path) as zf: xml=zf.read('word/document.xml')
        root=ET.fromstring(xml); ns={'w':'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
        return '\n'.join(''.join(t.text or '' for t in para.findall('.//w:t',ns)) for para in root.findall('.//w:p',ns))
    if suf=='.pdf':
        try:
            import fitz
            doc=fitz.open(path); return '\n'.join(page.get_text() for page in doc)
        except Exception:
            import PyPDF2
            with open(path,'rb') as f: return '\n'.join((page.extract_text() or '') for page in PyPDF2.PdfReader(f).pages)
    if suf=='.pptx':
        import zipfile
        from xml.etree import ElementTree as ET
        ns={'a':'http://schemas.openxmlformats.org/drawingml/2006/main'}; chunks=[]
        with zipfile.ZipFile(path) as zf:
            for name in sorted(n for n in zf.namelist() if n.startswith('ppt/slides/slide') and n.endswith('.xml')):
                root=ET.fromstring(zf.read(name)); texts=[t.text or '' for t in root.findall('.//a:t',ns)];
                if texts: chunks.append('\n'.join(texts))
        return '\n'.join(chunks)
    raise ValueError(f'Unsupported format: {path.suffix}')

def parse_blocks(text:str):
    lines=text.splitlines(); out={}; i=0
    while i < len(lines):
        m=SECTION_RE.match(lines[i].strip())
        if not m: i+=1; continue
        num=int(m.group(1)); path=m.group(2); i+=1; buf=[]
        while i < len(lines) and not END_RE.match(lines[i].strip()):
            buf.append(lines[i]); i+=1
        out.setdefault(num,{})[path]='\n'.join(buf).strip(); i+=1
    return out

def build_cs_dataset(en_json,translated_doc,out_json):
    en_data=load_json(en_json); grouped=parse_blocks(read_translation_text(Path(translated_doc))); out=[]
    for item in en_data:
        num=item['number']; tr=grouped.get(num,{})
        out.append({'number':num,'character':item.get('character',''),'traditional':item.get('traditional',''),'pinyin':item.get('pinyin',''),'name_cs':tr.get('name',''),'title_cs':tr.get('title',''),'above':{'chinese':get_path(item,'above.chinese',''),'symbolic_cs':tr.get('above.symbolic',''),'alchemical_cs':tr.get('above.alchemical','')},'below':{'chinese':get_path(item,'below.chinese',''),'symbolic_cs':tr.get('below.symbolic',''),'alchemical_cs':tr.get('below.alchemical','')},'judgment_cs':{'text':tr.get('judgment.text',''),'comments':tr.get('judgment.comments','')},'image_cs':{'text':tr.get('image.text',''),'comments':tr.get('image.comments','')},'lines_cs':{str(i):{'text':tr.get(f'lines.{i}.text',''),'comments':tr.get(f'lines.{i}.comments','')} for i in range(1,7)},'source_notes':{'english_name':item.get('name',''),'english_title':item.get('title',''),'english_judgment_text':get_path(item,'judgment.text',''),'english_image_text':get_path(item,'image.text','')}})
    Path(out_json).write_text(json.dumps(out,ensure_ascii=False,indent=2),encoding='utf-8')

def main():
    p=argparse.ArgumentParser(); sub=p.add_subparsers(dest='cmd',required=True)
    a=sub.add_parser('export-txt'); a.add_argument('--en-json',required=True); a.add_argument('--out-txt',required=True)
    b=sub.add_parser('export-html'); b.add_argument('--en-json',required=True); b.add_argument('--out-html',required=True)
    c=sub.add_parser('build-cs'); c.add_argument('--en-json',required=True); c.add_argument('--translated-doc',required=True); c.add_argument('--out-json',required=True)
    args=p.parse_args()
    if args.cmd=='export-txt': export_txt(args.en_json,args.out_txt)
    elif args.cmd=='export-html': export_html(args.en_json,args.out_html)
    else: build_cs_dataset(args.en_json,args.translated_doc,args.out_json)
if __name__=='__main__': main()
