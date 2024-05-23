// myCrypto is our own implementation of commonly used cryptographic functions, tailored to our project's needs.
// Cryptographic functions in this file are never implemented in house, but instead repurposed for our needs.
package myCrypto

import (
	"crypto/md5"
	"strconv"
)

// GetMD5Hash hashes an input string with MD5.
//
// Description:
// Generate a new MD5 hashed string based on the input string provided. If the provided string is unique, then the
// output hash will most likely also be unique. Please DO NOT use this algorithm for password hashing, as it is
// considered weak and prone to collisions.
//
// Parameters:
// - text: The returning hash is generated based on this string. To ensure a random string, you can use time.Now().
//
// Returns:
// random hash based on the current timestamp.
//
// Example:
// import time
// hash := myCrypto.GetMD5Hash(time.Now().String())
//
// Credit:
// https://stackoverflow.com/a/25286918
//
// Notes:
// The function outputs a string of numerical values. The reason is this function is currently only used for
// producing IDs used for registrations and notifications. In our implementations it is easier to retrieve ID
// if it is a numerical value, formatted as a string.
//
// Disclaimers:
// Do not use this function to hash passwords. MD5 is notoriously known for collisions.
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	nums := bytesToInts(hash)

	out := ""
	for i := 0; i < len(nums); i++ {
		out += strconv.Itoa(nums[i])
	}

	return out
}

// bytesToInts converts an array of bytes to a series of numerical integer values
// Currently this function is only used for producing IDs. This is far from optimal as the length of IDs
// vary, but it helps to retrieve IDs from URLs.
//
// Parameters:
// - byteArray: the array to convert to series of integers
//
// Returns:
// An array of integers.
func bytesToInts(byteArray [16]byte) []int {
	intSequence := make([]int, len(byteArray))
	for i, b := range byteArray {
		intSequence[i] = int(b)
	}
	return intSequence
}
