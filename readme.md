wanchunni@berkeley.edu

### Execute Steps

1. Run "main.go"

2. enter URL.  example url:
   `127.0.0.1:9527/?query="c1 == "abc" and * != "2" or c3 &= "123""`

3. After successful search, it will log "select succesfully! " in the console as well as the web page. The result dataset is saved in "result.csv", the original dataset is saved in "data.csv"
4. If failed, the error message will be shown on the page



### Design Consideration

1. Because the query can only be execute by URL and the original parser of the URL could treat special characters as different statement, using the format of `"?"` and treat the whole query as a string could avoid that kind of problem.
2. To parse the logical statement. I converted it into the format of arithmetic. *and* is multiply and *or* is plus. The the statement could be `num1 + num2 * num3 + num4 * num5 + num6` where the num is only 0 or 1. After calculating the value of the expression, it can be decided that the row in the dataset should be added to the result or not.
3. All the values in the dataset would be treated as string.



### Lessons Learnt

1. First time to write a http server in Golang
2. The way to parse simple logical statement
3. URL's input sanitization
4. The package management in Golang



