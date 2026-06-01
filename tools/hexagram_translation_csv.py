#!/usr/bin/env python3
import argparse
import csv
import json
from pathlib import Path

TEXT_FIELDS = [
    ("name", "name"),
    ("title", "title"),
    ("above.symbolic", "above.symbolic"),
    ("above.alchemical", "above.alchemical"),
    ("below.symbolic", "below.symbolic"),
    ("below.alchemical", "below.alchemical"),
    ("judgment.text", "judgment.text"),
    ("judgment.comments", "judgment.comments"),
    ("image.text", "image.text"),
    ("image.comments", "image.comments"),
]
for i in range(1, 7):
    TEXT_FIELDS.append((f"lines.{i}.text", f"lines.{i}.text"))
    TEXT_FIELDS.append((f"lines.{i}.comments", f"lines.{i}.comments"))

def get_path(obj, path, default=""):
    cur = obj
    for part in path.split('.'):
        if isinstance(cur, dict):
            cur = cur.get(part, default)
        else:
            return default
    return cur if cur is not None else default

def set_path(obj, path, value):
    parts = path.split('.')
    cur = obj
    for part in parts[:-1]:
        if part not in cur or not isinstance(cur[part], dict):
            cur[part] = {}
        cur = cur[part]
    cur[parts[-1]] = value

def export_csv(en_json: Path, out_csv: Path):
    data = json.loads(en_json.read_text(encoding='utf-8'))
    with out_csv.open('w', encoding='utf-8', newline='') as f:
        writer = csv.DictWriter(f, fieldnames=[
            'hexagram_number', 'character', 'traditional', 'pinyin',
            'section', 'source_path', 'source_text_en', 'target_text_cs'
        ])
        writer.writeheader()
        for item in data:
            for section, path in TEXT_FIELDS:
                writer.writerow({
                    'hexagram_number': item.get('number', ''),
                    'character': item.get('character', ''),
                    'traditional': item.get('traditional', ''),
                    'pinyin': item.get('pinyin', ''),
                    'section': section,
                    'source_path': path,
                    'source_text_en': get_path(item, path, ''),
                    'target_text_cs': ''
                })

def build_cs_dataset(en_json: Path, translated_csv: Path, out_json: Path):
    en_data = json.loads(en_json.read_text(encoding='utf-8'))
    grouped = {}
    with translated_csv.open('r', encoding='utf-8', newline='') as f:
        for row in csv.DictReader(f):
            num = int(row['hexagram_number'])
            grouped.setdefault(num, {})[row['source_path']] = row.get('target_text_cs', '').strip()
    out = []
    for item in en_data:
        num = item['number']
        tr = grouped.get(num, {})
        entry = {
            'number': num,
            'character': item.get('character', ''),
            'traditional': item.get('traditional', ''),
            'pinyin': item.get('pinyin', ''),
            'name_cs': tr.get('name', ''),
            'title_cs': tr.get('title', ''),
            'above': {
                'chinese': get_path(item, 'above.chinese', ''),
                'symbolic_cs': tr.get('above.symbolic', ''),
                'alchemical_cs': tr.get('above.alchemical', '')
            },
            'below': {
                'chinese': get_path(item, 'below.chinese', ''),
                'symbolic_cs': tr.get('below.symbolic', ''),
                'alchemical_cs': tr.get('below.alchemical', '')
            },
            'judgment_cs': {
                'text': tr.get('judgment.text', ''),
                'comments': tr.get('judgment.comments', '')
            },
            'image_cs': {
                'text': tr.get('image.text', ''),
                'comments': tr.get('image.comments', '')
            },
            'lines_cs': {
                str(i): {
                    'text': tr.get(f'lines.{i}.text', ''),
                    'comments': tr.get(f'lines.{i}.comments', '')
                } for i in range(1, 7)
            },
            'source_notes': {
                'english_name': item.get('name', ''),
                'english_title': item.get('title', ''),
                'english_judgment_text': get_path(item, 'judgment.text', ''),
                'english_image_text': get_path(item, 'image.text', '')
            }
        }
        out.append(entry)
    out_json.write_text(json.dumps(out, ensure_ascii=False, indent=2), encoding='utf-8')

def main():
    parser = argparse.ArgumentParser(description='Export EN hexagram dataset to CSV and build CS dataset from translated CSV.')
    sub = parser.add_subparsers(dest='cmd', required=True)

    p1 = sub.add_parser('export-csv')
    p1.add_argument('--en-json', required=True)
    p1.add_argument('--out-csv', required=True)

    p2 = sub.add_parser('build-cs')
    p2.add_argument('--en-json', required=True)
    p2.add_argument('--translated-csv', required=True)
    p2.add_argument('--out-json', required=True)

    args = parser.parse_args()
    if args.cmd == 'export-csv':
        export_csv(Path(args.en_json), Path(args.out_csv))
    elif args.cmd == 'build-cs':
        build_cs_dataset(Path(args.en_json), Path(args.translated_csv), Path(args.out_json))

if __name__ == '__main__':
    main()
