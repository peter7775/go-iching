# Překladový workflow bez CSV

Tento nástroj používá dokumentové formáty místo CSV.

Podporované vstupy pro překlad:
- `.txt`
- `.htm`
- `.html`
- `.docx`
- `.pdf`
- `.pptx`

Doporučené výstupní formáty pro export k překladu:
- `.txt`
- `.html`

## 1. Export EN datasetu do TXT

```bash
python3 tools/hexagram_translation_docs.py export-txt \
  --en-json dataembed/hexagrams.en.json \
  --out-txt output/hexagrams-to-translate.txt
```

## 2. Export EN datasetu do HTML

```bash
python3 tools/hexagram_translation_docs.py export-html \
  --en-json dataembed/hexagrams.en.json \
  --out-html output/hexagrams-to-translate.html
```

## 3. Generování českého datasetu z přeloženého dokumentu

```bash
python3 tools/hexagram_translation_docs.py build-cs \
  --en-json dataembed/hexagrams.en.json \
  --translated-doc output/hexagrams-to-translate.txt \
  --out-json dataembed/hexagrams.cs.json
```

Stejný příkaz funguje i pro `.html`, `.docx`, `.pdf` a `.pptx`, pokud dokument zachová značky:
- `@@HEX|...|START@@`
- `@@HEX|...|END@@`

## Důležité pravidlo

Překládej pouze text mezi START/END značkami a značky samotné nikdy neměň. Tyto značky slouží jako mapovací kotvy pro zpětné složení českého datasetu.

## Poznámky

- `.docx` je čitelný bez externího Word API, text se načítá přímo z OOXML.
- `.pdf` se čte přes `PyMuPDF` (`fitz`) nebo fallback `PyPDF2`.
- `.pptx` se čte přímo z XML uvnitř balíčku PowerPointu.
- Nejbezpečnější workflow je `TXT -> překlad -> build-cs` nebo `HTML -> překlad -> build-cs`.
