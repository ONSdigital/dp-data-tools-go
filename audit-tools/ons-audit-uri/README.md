# ons-audit-uri
1. Crawl whole of ONS site to check URI links that should respond do so.
2. Optionally save site info read to populate mongodb database for other research.

## read the instructions in the first 100 lines of comments in main.go

The notes will help you decide on how you want to run the app.

The default flag settings run the app in the most time optimal way.

With the default flags the app will output two reports of any
broken links in the `observations` directory:
```
broken_links.txt
broken_links_without_versions.txt
```
The URI links are searched from the home page down in a breadth before depth manner to have any broken links found listed in the above two files in the order in which a user might come across them.
