package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	// Don't be alerted by the following (API Call|Raw Data call) mismatches
	falseAlarms = map[string]bool{
		"AA|A": true,
		"CC|C": true,
		"GG|G": true,
		"TT|T": true,
		"DD|D": true,
		"II|I": true,
		"__|":  true,
		"--|":  true,
	}
	filenameRawdata string
	filenameAPIdata string
	filenameKey     string
)

// CallPair correlates an APICall (AA) with a RawDataCall (hopefully also AA)
type CallPair struct {
	APICall     string
	RawDataCall string
}

// Mismatch is the type of mismatch and how many times they occur in the genome
type Mismatch struct {
	CallPair
	Count int
}

// GenomesEndpoint is a container to unmarshal data from
// api.23andme.com/1/genomes/:profile_id/
type GenomesEndpoint struct {
	ID     string `json:"id"`
	Genome string
}

// SNP is just a more semantic name for something like rs124814
type SNP string

// Mismatches is a simple array of Mismatch
type Mismatches []Mismatch

// For sorting
func (m Mismatches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Mismatches) Len() int           { return len(m) }
func (m Mismatches) Less(i, j int) bool { return m[i].Count > m[j].Count }

func getSNPstoCall(filenameRawdata string) *map[string]string {
	var (
		file      *os.File
		lineBytes []byte
		err       error
	)
	if file, err = os.Open(filenameRawdata); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	SNPtoCall := make(map[string]string, 1050000)
	for {
		if lineBytes, _, err = reader.ReadLine(); err != nil {
			break
		}
		line := string(lineBytes)
		if strings.HasPrefix(line, "#") {
			continue
		}
		val := strings.Split(line, "\t")
		SNPtoCall[val[0]] = val[3]
	}
	return &SNPtoCall
}

func getIndexToSNP(filenameKey string) *map[int64]string {
	var (
		file      *os.File
		lineBytes []byte
		err       error
	)
	indexToSNP := make(map[int64]string, 1050000)
	if file, err = os.Open(filenameKey); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		if lineBytes, _, err = reader.ReadLine(); err != nil {
			break
		}
		line := string(lineBytes)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "index") {
			continue
		}
		val := strings.Split(line, "\t")
		var index int64
		if index, err = strconv.ParseInt(val[0], 10, 32); err != nil {
			break
		}
		indexToSNP[index] = val[1]
	}
	return &indexToSNP
}

func getCallpairs(filenameAPIdata string, indexToSNP *map[int64]string,
	SNPtoCall *map[string]string) (callpairs map[CallPair][]SNP, correct, incorrect int) {
	var err error
	callpairs = make(map[CallPair][]SNP, 10)
	jsondata, err := ioutil.ReadFile(filenameAPIdata)
	if err != nil {
		log.Fatal(err)
	}
	var genomes GenomesEndpoint
	json.Unmarshal(jsondata, &genomes)
	for index := 0; index < len(genomes.Genome); index += 2 {
		apiCall := fmt.Sprintf("%s%s", string(genomes.Genome[index]), string(genomes.Genome[index+1]))
		snpstr, _ := (*indexToSNP)[int64(index/2)]
		rawdataCall, _ := (*SNPtoCall)[snpstr]
		snp := SNP(snpstr)
		// Add mismatches; some are not true mismatches
		falseAlarm := falseAlarms[fmt.Sprintf("%s|%s", apiCall, rawdataCall)]
		if (apiCall != rawdataCall) && !falseAlarm {
			callpair := CallPair{APICall: apiCall, RawDataCall: rawdataCall}
			if _, found := callpairs[callpair]; !found {
				callpairs[callpair] = []SNP{snp}
			} else {
				callpairs[callpair] = append(callpairs[callpair], snp)
			}
			incorrect += 1
		} else {
			correct += 1
		}
	}
	return
}

func printAndCalculateMismatches(callpairs map[CallPair][]SNP, correct, incorrect int) {
	mismatches := Mismatches{}
	for callpair, snps := range callpairs {
		mismatch := Mismatch{CallPair: CallPair{APICall: callpair.APICall, RawDataCall: callpair.RawDataCall}, Count: len(snps)}
		mismatches = append(mismatches, mismatch)
	}
	sort.Sort(mismatches)
	for _, mismatch := range mismatches {
		fmt.Printf("APICall: %s\tRawDataCall: %s\tTotal: %d\t\n", mismatch.APICall, mismatch.RawDataCall, mismatch.Count)
		buffer := bytes.Buffer{}
		buffer.WriteString("SNPS: ")
		for i, snp := range callpairs[mismatch.CallPair] {
			buffer.WriteString(fmt.Sprintf("%s, ", snp))
			if (i%6 == 0) && (i > 0) {
				buffer.WriteString("\n")
			}
		}
		buffer.WriteString("\n\n")
		fmt.Print(buffer.String())
	}
	fmt.Printf("Same: %d, Mismatches: %d, Same: %f%%", correct, incorrect, float32(correct)/float32(incorrect+correct)*100)
}

func init() {
	flag.StringVar(&filenameRawdata, "r", "", "filename of raw data from https://www.23andme.com/you/download/ file (unzipped)")
	flag.StringVar(&filenameAPIdata, "a", "", "filename of API data from https://api.23andme.com/1/genomes/:profile_id/")
	flag.StringVar(&filenameKey, "k", "", "filename of downloaded https://api.23andme.com/res/txt/snps.data")
}

func main() {
	flag.Parse()
	// Require all command-line arguments
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == f.DefValue {
			fmt.Printf("Must pass -%s: %s\n", f.Name, f.Usage)
			os.Exit(1)
		}
	})

	SNPtoCall := getSNPstoCall(filenameRawdata)
	indexToSNP := getIndexToSNP(filenameKey)
	callpairs, correct, incorrect := getCallpairs(filenameAPIdata, indexToSNP, SNPtoCall)
	printAndCalculateMismatches(callpairs, correct, incorrect)
}
