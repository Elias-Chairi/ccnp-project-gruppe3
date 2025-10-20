// Package selectors provides node and actuator selection TLVs for COMMAND messages.
//
// Encoding rules:
//   - SINGLE_*: length=1, value=ID
//   - *_LIST:   length=N, value=ID1,ID2,...
//   - *_TYPE:   length=N, value=type name bytes
//   - ALL_*:    length=0, value=(empty)
package selectors
