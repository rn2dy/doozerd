package paxos

import (
    "fmt"
    "strings"
    "strconv"

    "borg/assert"
    "testing"
)

const (
    iSender = iota
    iCmd
    iRnd
    iNumParts
)

func accept(quorum int, ins, outs chan string) {
    var rnd, vrnd uint64
    var vval string

    ch, sent := make(chan int), 0
    for in := range ins {
        parts := strings.Split(in, ":", 3)
        i, _ := strconv.Btoui64(parts[iRnd], 10)
        switch {
            case i <= rnd:
            case i > rnd:
                rnd = i

                sent++
                msg := fmt.Sprintf("ACCEPT:%d:%d:%s", i, vrnd, vval)
                go func(msg string) { outs <- msg ; ch <- 1 }(msg)
        }
    }

    for x := 0; x < sent; x++ {
        <-ch
    }

    close(outs)
}

func TestIgnoresStaleInvites(t *testing.T) {
    ins := make(chan string)
    outs := make(chan string)

    exp := "ACCEPT:2:0:"

    go accept(2, ins, outs)
    // Send a message with no senderId
    ins <- "1:INVITE:2"
    ins <- "1:INVITE:1"
    close(ins)

    got := ""
    for x := range outs {
        got += x
    }

    // outs was closed; therefore all messages have been processed
    assert.Equal(t, exp, got, "")
}