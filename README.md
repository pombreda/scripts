scripts
=======

You need a go compiler (http://golang.org/doc/install) for this to work.

You need to download three files: the index-to-SNP key from the API,
the JSON response of the call to api.23andme.com/1/genomes/, and the
raw 23andMe data from 23andme.com. Details in the script.

When you have these three files in your directory, run the script. It will
output the mismatches between the API and the raw 23andMe download.

usage:

    go run compare_api_raw_download.go -a apidata.txt -k snps.data -r rawdata.txt

output:

    (...)
    ApiCall: AT     RawDataCall:    Total: 4
    SNPS: rs2557018, rs2571902, rs2573893, rs2751692,

    ApiCall: CG     RawDataCall:    Total: 3
    SNPS: rs2535092, rs2557825, rs2915713,

    ApiCall: TT     RawDataCall:    Total: 2
    SNPS: rs1152098, rs35844236,

    ApiCall: DI     RawDataCall:    Total: 1
    SNPS: i4000257,

    ApiCall: __     RawDataCall: TT Total: 1
    SNPS: rs429358,

    ApiCall: DI     RawDataCall: D  Total: 1
    SNPS: rs3838646,

    Same: 1046257, Mismatches: 1701, Same: 99.837685%

Note that the vast majority of "mismatches" are false alarms, that
come from hemizygous calls that have been smushed in the RawData,
but not in the API, (i.e., ApiCall: BB, RawDataCall: B).
They are not included in the output, but you may change that with the variable
```FALSE_ALARMS``` in the script.
