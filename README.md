# Tokens

### MongoDB Compass:
	mongodb+srv://medods:yD0D3vEXWCtnFkvZ@cluster0.ufzsk.mongodb.net/test
##

### Генерация токенов:
Необходима сделать POST запрос на url: http://tokens-acc-ref.herokuapp.com/generate

Входные данные:

	{
	   "user_id" 64842
	   "GUID"   "FF4vjiQAU8e6KYq/edxsrQ=="
	}

Ответ:

	{
	   "access": "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJ1c2VyX2lkIjowLAiOjE2MDIwMDQwMjcsIk.MaHBfL0C7lsr7i5joFQ5VT3NdPDxCHKQ9ucFnPw/xaiRNOuFHvE4ksWDd6tdIOmKyQNYTeW2c06oPYgpJXqteg==",
	   "refresh": "eyJhY2Nlc3MiOiJNdGNJV3I3MD9IiwiZXhCJHVUlEIjoiQ2IxQUtDWkZad0puMG1tQ1ROTHg2QT09In0=.1d1IbXdFcEHZA1NC3C8SeZOEH62oraFbbhKK2SlEw63v5IOZIqblobWNJSc9OYnJ1YKRj3ndhXa87o3bFbbEAw=="
	}
##
### Операция refresh над токенами:
Необходима сделать POST запрос на url: http://tokens-acc-ref.herokuapp.com/refresh

Входные данные:

	{
	   "access": "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJ1c2VyX2lkIjowLAiOjE2MDIwMDQwMjcsIk.MaHBfL0C7lsr7i5joFQ5VT3NdPDxCHKQ9ucFnPw/xaiRNOuFHvE4ksWDd6tdIOmKyQNYTeW2c06oPYgpJXqteg==",
  	   "refresh": "eyJhY2Nlc3MiOiJNdGNJV3I3MD9IiwiZXhCJHVUlEIjoiQ2IxQUtDWkZad0puMG1tQ1ROTHg2QT09In0=.1d1IbXdFcEHZA1NC3C8SeZOEH62oraFbbhKK2SlEw63v5IOZIqblobWNJSc9OYnJ1YKRj3ndhXa87o3bFbbEAw=="
	}

Ответ:

	{
	   "access": "eyJhbGciOiJTSEE1MTIiLCJ0eXAiOiJKV1QifQ==.eyJ1c2VyX2lkIjowLAiOjE2MDIwMDQwMjcsIk.MaHBfL0C7lsr7i5joFQ5VT3NdPDxCHKQ9ucFnPw/xaiRNOuFHvE4ksWDd6tdIOmKyQNYTeW2c06oPYgpJXqteg==",
	   "refresh": "eyJhY2Nlc3MiOiJNdGNJV3I3MD9IiwiZXhCJHVUlEIjoiQ2IxQUtDWkZad0puMG1tQ1ROTHg2QT09In0=.1d1IbXdFcEHZA1NC3C8SeZOEH62oraFbbhKK2SlEw63v5IOZIqblobWNJSc9OYnJ1YKRj3ndhXa87o3bFbbEAw=="
	}
##

### Удаление Refresh token:
Необходима сделать POST запрос на url: http://tokens-acc-ref.herokuapp.com/delete

Входные данные:

	{
	   "refresh": "eyJhY2Nlc3MiOiJNdGNJV3I3MD9IiwiZXhCJHVUlEIjoiQ2IxQUtDWkZad0puMG1tQ1ROTHg2QT09In0=.1d1IbXdFcEHZA1NC3C8SeZOEH62oraFbbhKK2SlEw63v5IOZIqblobWNJSc9OYnJ1YKRj3ndhXa87o3bFbbEAw=="
	}

Ответ:

    http.StatusOK
##

### Удаление Refresh tokens по id:
Необходима сделать POST запрос на url: http://tokens-acc-ref.herokuapp.com/deleteall

Входные данные:

	{
	   "user_id": 64842
	}

Ответ:

    http.StatusOK
