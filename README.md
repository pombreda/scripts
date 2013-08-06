scripts
=======

Collection of helper scripts

You need a go compiler (http://golang.org/doc/install) for this to work.

You need to download three files: the index-to-SNP key from the API,
the JSON response of the call to api.23andme.com/1/genomes/, and the
raw 23andMe data from 23andme.com. Details in the script.

When you have these three files in your directory, run the script. It will
output the mismatches between the API and the raw 23andMe download.

usage:

    go run compare_api_raw_download.go

output:

    2013/08/05 19:41:04 ApiCall: AA RawDataCall:    Total: 26
    SNPS: rs1006094, rs10215320, rs10231034, rs10487532, rs11535222, rs12533751,
    rs12534123,
    rs13223129, rs13233002, rs16622, rs17163497, rs17708955, rs17767132,
    rs2269332, rs2665033, rs2734189, rs2855951, rs2855954, rs2855983,
    rs361379, rs7210806, rs7501783, rs7503902, rs7787291, rs8066263,
    rs8080666,

    2013/08/05 19:41:04 ApiCall: GG RawDataCall:    Total: 20
    SNPS: rs2855945, rs12603312, rs12667623, rs12672113, rs17163396, rs17163403,
    rs17257,
    rs17282, rs17286, rs2040351, rs2213198, rs2213199, rs2252241,
    rs2734151, rs2855150, rs2855943, rs361473, rs7405659, rs8065316,
    rs8069746,

    2013/08/05 19:41:04 ApiCall: CC RawDataCall:    Total: 7
    SNPS: rs1006093, rs2040350, rs6947359, rs6968260, rs6979080, rs7789642,
    rs8071835,


    2013/08/05 19:41:04 ApiCall: AG RawDataCall:    Total: 4
    SNPS: rs11079538, rs2070782, rs4968624, rs4968723,

    2013/08/05 19:41:04 ApiCall: CT RawDataCall: -- Total: 1
    SNPS: rs5929791,

    2013/08/05 19:41:04 Same: 523921, Mismatches: 58, Same: 99.988930%

Note that the vast majority of "mismatches" are false alarms, that
come from hemizygous calls that have been smushed in the RawData,
but not in the API, (i.e., ApiCall: BB, RawDataCall: B).
They are not included in the output, but you may change that with the variable
```FALSE_ALARMS``` in the script.
