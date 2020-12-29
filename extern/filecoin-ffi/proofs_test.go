package ffi

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/filecoin-project/filecoin-ffi/generated"

	"github.com/stretchr/testify/assert"

	commcid "github.com/filecoin-project/go-fil-commcid"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	fmt.Println("hello")
}

func TestRegisteredSealProofFunctions(t *testing.T) {
	fmt.Println("===hello")
	WorkflowRegisteredSealProofFunctions(newTestingTeeHelper(t))
}


func TestRegisteredPoStProofFunctions(t *testing.T) {
	WorkflowRegisteredPoStProofFunctions(newTestingTeeHelper(t))
}

/*
/**
	官方注释：Seek设置下一次读/写的位置。offset为相对偏移量，而whence决定相对位置：0为相对文件开头，1为相对当前位置，2为相对文件结尾。它返回新的偏移量（相对开头）和可能的错误。
	whence参数
		io.SeekStart 	// 0
		io.SeekCurrent 	// 1
		io.SeekEnd 		// 2
*/

func TestReadFile(t *testing.T) {
	f, err := os.OpenFile(`a.txt`, os.O_RDWR, os.ModePerm)
	fmt.Println("err", err)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	end, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("end: ", end) // end： 10
	fs, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("size: ", fs.Size())

	// 都是含前不含后的概念
	// offset是从0开始的, 可以比当前的文件内容长度大，多出的部分会用空(0)来代替
	start, err := f.Seek(13, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("start: ", start)
	_, err = f.WriteString("a")
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 102)
	n, err := f.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b[:n]) // [48 49 50 51 52 53 54 55 56 57 0 0 97]
	fmt.Println(string(b[:n])) //[35 52 53 10 32 32 32 32 104 101 108 108 111 32 119 111 114 108 100]
}

func TestProofsLifecycle(t *testing.T) {
	WorkflowProofsLifecycle(newTestingTeeHelper(t))
}

func TestGetGPUDevicesDoesNotProduceAnError(t *testing.T) {
	fmt.Println("==hello==")
	WorkflowGetGPUDevicesDoesNotProduceAnError(newTestingTeeHelper(t))
}

func TestGenerateWinningPoStSectorChallenge(t *testing.T) {
	WorkflowGenerateWinningPoStSectorChallenge(newTestingTeeHelper(t))
}

func TestGenerateWinningPoStSectorChallengeEdgeCase(t *testing.T) {
	WorkflowGenerateWinningPoStSectorChallengeEdgeCase(newTestingTeeHelper(t))
}

func TestJsonMarshalSymmetry(t *testing.T) {
	for i := 0; i < 100; i++ {
		xs := make([]publicSectorInfo, 10)
		for j := 0; j < 10; j++ {
			var x publicSectorInfo
			var commR [32]byte
			_, err := io.ReadFull(rand.Reader, commR[:])
			require.NoError(t, err)

			// commR is defined as 32 long above, error can be safely ignored
			x.SealedCID, _ = commcid.ReplicaCommitmentV1ToCID(commR[:])

			n, err := rand.Int(rand.Reader, big.NewInt(500))
			require.NoError(t, err)
			x.SectorNum = abi.SectorNumber(n.Uint64())
			xs[j] = x
		}
		toSerialize := newSortedPublicSectorInfo(xs...)

		serialized, err := toSerialize.MarshalJSON()
		require.NoError(t, err)

		var fromSerialized SortedPublicSectorInfo
		err = fromSerialized.UnmarshalJSON(serialized)
		require.NoError(t, err)

		require.Equal(t, toSerialize, fromSerialized)
	}
}

func TestDoesNotExhaustFileDescriptors(t *testing.T) {
	m := 500         // loops
	n := uint64(508) // quantity of piece bytes

	for i := 0; i < m; i++ {
		// create a temporary file over which we'll compute CommP
		file, err := ioutil.TempFile("", "")
		if err != nil {
			panic(err)
		}

		// create a slice of random bytes (represents our piece)
		b := make([]byte, n)

		// load up our byte slice with random bytes
		if _, err = rand.Read(b); err != nil {
			panic(err)
		}

		// write buffer to temp file
		if _, err := bytes.NewBuffer(b).WriteTo(file); err != nil {
			panic(err)
		}

		// seek to beginning of file
		if _, err := file.Seek(0, 0); err != nil {
			panic(err)
		}

		if _, err = GeneratePieceCID(abi.RegisteredSealProof_StackedDrg2KiBV1, file.Name(), abi.UnpaddedPieceSize(n)); err != nil {
			panic(err)
		}

		if err = file.Close(); err != nil {
			panic(err)
		}
	}
}

func newTestingTeeHelper(t *testing.T) *testingTeeHelper {
	return &testingTeeHelper{t: t}
}

type testingTeeHelper struct {
	t *testing.T
}

func (tth *testingTeeHelper) RequireTrue(value bool, msgAndArgs ...interface{}) {
	require.True(tth.t, value, msgAndArgs)
}

func (tth *testingTeeHelper) RequireNoError(err error, msgAndArgs ...interface{}) {
	require.NoError(tth.t, err, msgAndArgs)
}

func (tth *testingTeeHelper) RequireEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	require.Equal(tth.t, expected, actual, msgAndArgs)
}

func (tth *testingTeeHelper) AssertNoError(err error, msgAndArgs ...interface{}) bool {
	return assert.NoError(tth.t, err, msgAndArgs)
}

func (tth *testingTeeHelper) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Equal(tth.t, expected, actual, msgAndArgs)
}

func (tth *testingTeeHelper) AssertTrue(value bool, msgAndArgs ...interface{}) bool {
	return assert.True(tth.t, value, msgAndArgs)
}

func TestProofTypes(t *testing.T) {
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWinning2KiBV1, abi.RegisteredPoStProof_StackedDrgWinning2KiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWinning8MiBV1, abi.RegisteredPoStProof_StackedDrgWinning8MiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWinning512MiBV1, abi.RegisteredPoStProof_StackedDrgWinning512MiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWinning32GiBV1, abi.RegisteredPoStProof_StackedDrgWinning32GiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWinning64GiBV1, abi.RegisteredPoStProof_StackedDrgWinning64GiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWindow2KiBV1, abi.RegisteredPoStProof_StackedDrgWindow2KiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWindow8MiBV1, abi.RegisteredPoStProof_StackedDrgWindow8MiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWindow512MiBV1, abi.RegisteredPoStProof_StackedDrgWindow512MiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWindow32GiBV1, abi.RegisteredPoStProof_StackedDrgWindow32GiBV1)
	assert.EqualValues(t, generated.FilRegisteredPoStProofStackedDrgWindow64GiBV1, abi.RegisteredPoStProof_StackedDrgWindow64GiBV1)

	assert.EqualValues(t, generated.FilRegisteredSealProofStackedDrg2KiBV1, abi.RegisteredSealProof_StackedDrg2KiBV1)
	assert.EqualValues(t, generated.FilRegisteredSealProofStackedDrg8MiBV1, abi.RegisteredSealProof_StackedDrg8MiBV1)
	assert.EqualValues(t, generated.FilRegisteredSealProofStackedDrg512MiBV1, abi.RegisteredSealProof_StackedDrg512MiBV1)
	assert.EqualValues(t, generated.FilRegisteredSealProofStackedDrg32GiBV1, abi.RegisteredSealProof_StackedDrg32GiBV1)
	assert.EqualValues(t, generated.FilRegisteredSealProofStackedDrg64GiBV1, abi.RegisteredSealProof_StackedDrg64GiBV1)
}
