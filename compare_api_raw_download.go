package main

import (
	"bufio"
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// You need these files present to run the check.
const (
	// Save https://www.23andme.com/you/download/
	// Unzip it, and rename genome_Your_Name_Full_timestamp.txt -> rawdata.txt
	RAWDATA_FILENAME = "rawdata.txt"
	// curl https://api.23andme.com/1/genomes/192889f1/ > apidata.txt
	API_DATA_FILENAME = "apidata.txt"
	// curl https://api.23andme.com/res/txt/snps.data
	KEY_FILENAME = "snps.data"
)

type CallPair struct {
	ApiCall     string
	RawDataCall string
}

type Mismatch struct {
	CallPair
	Count int
}

type GenomesEndpoint struct {
	Id     string
	Genome string
}

type SNP string

type Mismatches []Mismatch

// For sorting
func (m Mismatches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Mismatches) Len() int           { return len(m) }
func (m Mismatches) Less(i, j int) bool { return m[i].Count > m[j].Count }

func getSNPstoCall() *map[string]string {
	var (
		file *os.File
		line []byte
		err  error
	)
	if file, err = os.Open(RAWDATA_FILENAME); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	SNPtoCall := make(map[string]string, 1050000)
	for {
		if line, _, err = reader.ReadLine(); err != nil {
			break
		}
		linestring := string(line)
		if strings.HasPrefix(linestring, "#") {
			continue
		}
		val := strings.Split(linestring, "\t")
		SNPtoCall[val[0]] = val[3]
	}
	return &SNPtoCall
}

func getIndexToSNP() *map[int64]string {
	var (
		file *os.File
		line []byte
		err  error
	)
	indexToSNP := make(map[int64]string, 1050000)
	if file, err = os.Open(KEY_FILENAME); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		if line, _, err = reader.ReadLine(); err != nil {
			break
		}
		linestring := string(line)
		if strings.HasPrefix(linestring, "#") || strings.HasPrefix(linestring, "index") {
			continue
		}
		val := strings.Split(linestring, "\t")
		var index int64
		if index, err = strconv.ParseInt(val[0], 10, 32); err != nil {
			break
		}
		indexToSNP[index] = val[1]
	}
	return &indexToSNP
}

func getCallpairs(indexToSNP *map[int64]string,
	SNPtoCall *map[string]string) (callpairs map[CallPair][]SNP, correct, incorrect int) {
	var err error
	callpairs = make(map[CallPair][]SNP, 10)
	jsondata, err := ioutil.ReadFile(API_DATA_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	var genomes GenomesEndpoint
	json.Unmarshal(jsondata, &genomes)
	for index := 0; index < len(genomes.Genome)/2-1; index += 2 {
		api_call := fmt.Sprintf("%s%s", string(genomes.Genome[index]), string(genomes.Genome[index+1]))
		snpstr, _ := (*indexToSNP)[int64(index/2)]
		raw_data_call, _ := (*SNPtoCall)[snpstr]
		snp := SNP(snpstr)
		// Add mismatches
		if api_call != raw_data_call {
			callpair := CallPair{ApiCall: api_call, RawDataCall: raw_data_call}
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
		mismatch := Mismatch{CallPair: CallPair{ApiCall: callpair.ApiCall, RawDataCall: callpair.RawDataCall}, Count: len(snps)}
		mismatches = append(mismatches, mismatch)
	}
	sort.Sort(mismatches)
	for _, mismatch := range mismatches {
		log.Printf("ApiCall: %s\tRawDataCall: %s\tTotal: %d\t\n", mismatch.ApiCall, mismatch.RawDataCall, mismatch.Count)
		//if mismatch.Count < 10 {
		//buffer := bytes.Buffer{}
		//buffer.WriteString("SNPS: ")
		//for i, snp := range callpairs[mismatch.CallPair] {
		//buffer.WriteString(fmt.Sprintf("%s, ", snp))
		//if (i%6 == 0) && (i > 0) {
		//buffer.WriteString("\n")
		//}
		//}
		//buffer.WriteString("\n\n")
		//fmt.Print(buffer.String())
		//}
	}
	log.Printf("Same: %d, Mismatches: %d, Same: %f%%", correct, incorrect, float32(correct)/float32(incorrect+correct)*100)
}

func main() {
	SNPtoCall := getSNPstoCall()
	indexToSNP := getIndexToSNP()
	callpairs, correct, incorrect := getCallpairs(indexToSNP, SNPtoCall)
	printAndCalculateMismatches(callpairs, correct, incorrect)
}
