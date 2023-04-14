# gsscrape

__gsscrape__ is utility for scraping data from HTML structurs and saving to output file. __gsscrape__ work in multithread&concurancy mode. You may set number of the workers working in same time in the command line flag. The data that will be scraped from HTML should specified in the code of _org.go_ file. The object __OrgHtmlJson__ imlementing __Scrape__ func of interface __Scraper__ is customize the HTML data processing and the output data format.

Usage:
---
__gsscrape <-h> <-t NNN> <-w NNN> <-o output_file> -i input_file <url>__
##### Flags:
-h -help:       Show help (Default: false)\
-t:             The timeout in seconds for waiting a responses from web sites. (Default: 5)\
-v -verbouse:   Output fool log to StdOut (Default: false)\
-w:             The number of workers working in the same time. (Default: 5)\
-o:             File for result output. If the flag is absent then output will to the StdOut.\
-i:             Input web src for scraping data. If the flag is absent then input should from last argument.

The list of URLs for processing defined in the input file or in the command line (one URL). The parameters in the URL can include the masks of the types:\
[nnn:nnn] - range between numbers not including last number (GO slice style)\
[word1;word2;word3] - enumeration of strings

For exaple the mask
```
html://www.site.com/path?chapter=[one;two]&page=[1:3]
```
will trasform to tje URLs:
```
html://www.site.com/path?chapter=one&page=1
html://www.site.com/path?chapter=one&page=2
html://www.site.com/path?chapter=two&page=1
html://www.site.com/path?chapter=two&page=2
```
The line in the input file may be commented by "//"
