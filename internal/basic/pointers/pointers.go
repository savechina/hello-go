package pointers

import (
	"fmt"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "pointers", Run)
}

func main() {
	Run()
}

type wallet struct {
	balance int
}

func (w *wallet) Deposit(amount int) {
	if w == nil {
		return
	}
	w.balance += amount
}

func (w *wallet) Balance() int {
	if w == nil {
		return 0
	}
	return w.balance
}

type learner struct {
	name string
}

func renameWithPointer(target *string, next string) bool {
	if target == nil {
		return false
	}
	*target = next
	return true
}

func swapValues(left *int, right *int) bool {
	if left == nil || right == nil {
		return false
	}
	*left, *right = *right, *left
	return true
}

func safeLearnerName(item *learner) string {
	if item == nil {
		return "nil learner"
	}
	return item.name
}

func exampleAddressAndDereference() string {
	value := 10
	pointer := &value
	*pointer += 5
	return fmt.Sprintf("value=%d pointer-set=%d", value, *pointer)
}

func examplePointerReceiver() string {
	account := &wallet{}
	account.Deposit(30)
	account.Deposit(12)
	return fmt.Sprintf("wallet balance=%d", account.Balance())
}

func exampleNilPointer() string {
	var nobody *learner
	var broken *wallet
	return fmt.Sprintf("learner=%s balance=%d", safeLearnerName(nobody), broken.Balance())
}

// Run prints the pointers chapter examples.
func Run() {
	fmt.Println("[pointers] example 1:", exampleAddressAndDereference())
	fmt.Println("[pointers] example 2:", examplePointerReceiver())
	fmt.Println("[pointers] example 3:", exampleNilPointer())
}
