package main

import (
	"fmt"
	"os"
	"time"
	"bufio"
	"strings"
	"sort"
	"hash/crc32"
	"strconv"
)

const FILE string	= "C:/Users/crisping/go/src/github.com/user/mining/project/data/sentences/sentences.txt"
const MIN_SEN int	= 10 	// min sentence length
var num_sentences int	//number of sentences in the file
var start time.Time

func main() {
	find_candidate_matches()
	fmt.Println("out of find_candidate_matches  ", time.Since(start))
	compare_candidates()
	print_results()
}

func find_candidate_matches() {
	start = time.Now()
	// read file
	file := openFile(FILE)
	defer file.Close()
	reader := bufio.NewReader(file)
//	fields := make([]string, 100)
	for lines := 1; ; lines++ {
		num_sentences += 1
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fields := strings.Fields(line)
		id, _ := strconv.Atoi(fields[0])
		sentence_lengths[id] = len(fields) - 1  	// fields[0] has sentence id
		createHashes(fields[1:], id)
	}
//	fmt.Println(sentence_lengths)
	fmt.Printf("Initial read and hashes created- lines: %v  time: %v\n", num_sentences, time.Since(start))
	hash_matches := 0
	for _, hasharr := range Hashes {
		if len(hasharr) > 1 {
			hash_matches += len(hasharr)
			filter_hashes_by_sentence_length(hasharr)
		}
	}
//	fmt.Println("hashes: ", Hashes)
	fmt.Printf("Number of hash matches: %v\n", hash_matches)
	fmt.Printf("Number of candidates for pass 2: %v\n", idsToCompareCount)
}

var Keep = make(map[int]int)

type CompareKey struct {
	id1, id2 int
}	// store sentence on pass2 map[sid]# of accesses  
var Matched = make(map[CompareKey]bool, 0)	// ids to compare on pass 2, higher id first

var idsToCompare = make(map[int][]int)	// will be used in pass2
var idsToCompareCount = 0

func addids(id1, id2 int) {
	if Matched[CompareKey{id1, id2}] == false {		// these sids are already flagged for comparison
		Matched[CompareKey{id1, id2}] = true
		Keep[id2] += 1 	// count of compares by higher sids
		idsToCompare[id1] = append(idsToCompare[id1], id2)
		idsToCompareCount += 1
	}	
}

func filter_hashes_by_sentence_length(a []int) {
	for i := 0; i < len(a); i++ {
		sli := sentence_lengths[a[i]]
		for j := i + 1; j < len(a); j++ {
			slj := sentence_lengths[a[j]]
//			fmt.Println(a[i], a[j], sentence_lengths[a[i]], sentence_lengths[a[j]])
			if sli  == slj || sli == slj -1 ||
				sli == slj + 1 {
				if a[i] < a[j] {
					addids(a[j], a[i])
				} else {
					addids(a[i], a[j])
				}
			}
		}
	}
}

var sentence_lengths = make(map[int]int, num_sentences)		// word count for each sentence id
var Hashes = make(map[uint32][]int, num_sentences*2)		// two hashes per sentence 

func createHashes(fields []string, id int) {
	sort.Strings(fields)
//	id, _ := strconv.Atoi(fields[0])
	hash1 := getHashSum(fields[:5])
	hash2 := getHashSum(fields[5:10])
	arr := []int{}
	arr = append(Hashes[hash1], id)
	Hashes[hash1] = arr
	arr = append(Hashes[hash2], id)
	Hashes[hash2] = arr
}

func getHashSum(str []string) uint32 {
	s := strings.Join(str, "")
	b := make([]byte, 0)
	dst := strconv.AppendQuote(b, s)
	h := crc32.NewIEEE()
	h.Write(dst)
	return h.Sum32()
}


// create word index by increasing frequency

func openFile(f string) *os.File {
	file, err := os.Open(FILE)
	if err != nil {
		fmt.Println(err)
	}
	return file
}

func compare_candidates() {
	return
}

func print_results() {
	return
}

