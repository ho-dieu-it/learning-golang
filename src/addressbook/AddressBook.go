package main

import (
	"fmt"
	"strings"
	"os"
	"encoding/json"
	"bufio"
	"io/ioutil"
	"strconv"
)

const ADDRESS_BOOK_FILE string = "AddressBook.json"

var choices = []string{
	"1. Add a contact",
	"2. List",
	"3. Find",
	"4. Delete",
	"Q. Quit",
}

type Contact struct {
	ID                       int
	Name, Email, PhoneNumber string
}

var contacts = []Contact{}

func main() {
	loadContacts()

	for {
		showMenu()

		switch getChoice() {
		case "1","2","3":
			contact := showAddContact()
			addContact(contact)
		case "2":
			showContactList(contacts)
		case "3":
			name := prompt("Enter a name:")
			result := findContact(name)
			showContactList(result)
		case "4":
			id := prompt("Enter contact ID:")
			deleteContact(id)
		case "q":
			os.Exit(0)
		}
	}
}

func showMenu() {
	line := "===================="

	fmt.Println(line)
	fmt.Println("MENU")
	fmt.Println(line)
	for _, choice := range choices {
		fmt.Println(choice)
	}
	fmt.Println(line)
}

func getChoice() string {
	fmt.Print("Enter your choice:")

	var choice string
	fmt.Scanln(&choice)

	if !isValidChoice(choice) {
		fmt.Println("ERROR! Invalid choice.")
		getChoice()
	}

	return choice
}

func isValidChoice(key string) bool {
	for _, choice := range choices {
		if key == strings.ToLower(string(choice[0])) {
			return true
		}
	}

	return false
}

func showAddContact() *Contact {
	contact := &Contact{}
	fmt.Println("Enter new contact")

	contact.Name = prompt("Name:");
	contact.Email = prompt("Email:");
	contact.PhoneNumber = prompt("Phone Number:");

	return contact
}

func prompt(text string) string {
	fmt.Print(text)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func showContactList(contacts []Contact) {
	fmt.Printf("CONTACTS (%d)\n", len(contacts))

	format := "| %2s | %-20s | %-30s | %-15s |\n"

	line := fmt.Sprintf(strings.Replace(format, " ", "", -1),
		strings.Repeat("-", 4), strings.Repeat("-", 22), strings.Repeat("-", 32), strings.Repeat("-", 17))

	fmt.Print(line)
	fmt.Printf(format, "ID", "Name", "Email", "Phone Number")
	fmt.Print(line)

	for _, contact := range contacts {
		fmt.Printf(format, strconv.Itoa(contact.ID), contact.Name, contact.Email, contact.PhoneNumber)
	}

	fmt.Print(line)
}

func addContact(contact *Contact) {
	newID := 1;

	if numberOfContacts := len(contacts); numberOfContacts > 0 {
		newID = contacts[numberOfContacts - 1].ID + 1
	}
	contact.ID = newID
	contacts = append(contacts, *contact)

	saveContacts()
}

func deleteContact(id string) {
	for i, contact := range contacts {
		if strconv.Itoa(contact.ID) != id {
			continue
		}

		contacts = append(contacts[:i], contacts[i + 1:]...)
		saveContacts()
		break
	}
}

func findContact(name string) []Contact {
	result := []Contact{}
	for _, contact := range contacts {
		if strings.Contains(contact.Name, name) {
			result = append(result, contact)
		}
	}
	return result
}

func saveContacts()  {
	json, _ := json.Marshal(contacts)
	file, _ := os.Create(ADDRESS_BOOK_FILE)
	file.Write(json)

	defer file.Close()
}

func loadContacts() {
	data, err := ioutil.ReadFile(ADDRESS_BOOK_FILE)
	if err != nil {
		return
	}
	json.Unmarshal([]byte(data), &contacts)
}