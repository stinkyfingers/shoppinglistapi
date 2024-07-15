User
ID uuid
FirstName string
LastName string
Email string
PictureURL string

Collection
Creator User 
Shared []User
Name string
Lists []List

List
Name string
Items []ListItem
Store Store

ListItem
Item Item
Quantity int
Status string
Location string

Item
ID uuid
Name string

Store
ID uuid
Name string
Location string

