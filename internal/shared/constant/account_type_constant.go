package constant

type AccountTypeEnum string

const (
	AccountTypeCustomer AccountTypeEnum = "user"
	AccountTypeAdmin    AccountTypeEnum = "admin"
	AccountTypeStaff    AccountTypeEnum = "staff"
)
