// thanks for the support of  r/u/nsd433
package main

import (
        "fmt"
)

type CtgTree struct {
        root *node
}
type node struct {
        Parent   *node
        Name     string
        ID       int
        Children []*node
}

func (c *CtgTree) Serialize() (buff []int) {
        ch := c.Walker()
        buff = make([]int, 0)
        for elem := range ch {
                fmt.Printf(" %d\n", elem)
                buff = append(buff, elem)
        }
        return buff
}

func (c *CtgTree) Walker() <-chan int {
        ch := make(chan int)
        go func() {
                walk(c.root, ch)
                close(ch)
        }()
        return ch
}

// Walk traverses a tree
// sending each Value on a channel.
// 0 is "(" represents nil node at the serialization
func walk(t *node, ch chan int) {
        ch <- t.ID
        for _, child := range t.Children {
                walk(child, ch)
        }
        ch <- 0
}

// Deserialize constructs the c from slice of int
func (c *CtgTree) Deserialize(s []int) (err error) {
        ch := make(chan int)
        go func() {
                for _, e := range s {
                        ch <- e
                }
                close(ch)
        }()

        c.root, err = build(nil, ch)
        return err
}

// build gets a nil node,
// and incoming int channel ch
// returns root of new tree
func build(t *node, ch chan int) (*node, error) {
        id := <-ch
        if t == nil {
                t = &node{ID: id}
                return build(t, ch)
        }

        if id == 0 {
                if t.Parent == nil {
                        for _ = range ch {
                                return nil, fmt.Errorf("invalid input sequence at s")
                        }
                        return t, nil
                }
                return build(t.Parent, ch)
        }
        n := &node{ID: id, Parent: t}
        t.Children = append(t.Children, n)
        return build(n, ch)
}

func main() {
        ctg := CtgTree{}
        buff := []int{1, 2, 4, 0, 5, 0, 6, 0, 0, 6, 0, 0, 3, 7, 0, 0, 0}
        fmt.Println("deserializing")
        err := ctg.Deserialize(buff)
        if err != nil {
                panic(err)
        }
        buff = append(buff)
        _ = ctg.Serialize()
}
