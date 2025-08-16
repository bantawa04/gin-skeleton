package constant

type AccountTypeEnum string

const (
	AccountTypeCustomer AccountTypeEnum = "customer"
	AccountTypeAdmin    AccountTypeEnum = "admin"
	AccountTypeStaff    AccountTypeEnum = "staff"
)