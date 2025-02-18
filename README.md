## About
This is a version 2 of my original HDB Project. The original project was just a basic interface to do query calls on the db and do basic averaging.

With this project, I have added the feature of using different models to predict future prices. Also, the feature of sorting HDB price data by area has been introduced, and the different areas will be compared.

![image](https://github.com/user-attachments/assets/78b2a866-2e18-4afc-b957-bd741b25e687)
![Screenshot 2025-01-12 at 3 52 53 PM](https://github.com/user-attachments/assets/1fe2913a-5a0e-4706-8a78-c9393e7fd2d8)



## How to Use
#### Running the server
1. clone the project using `git clone https://github.com/bedminer1/HDB2.git`
2. go into the api directory using the following command `cd backend/cmd/api`
3. run the server using the following command `go run .`

The server is now running, here is a list of endpoints.


### Fetching 
1. /records gets you a list of HDB records 
2. /monthly_stats aggregates HDB records into a record representing the whole month, then showing records for different months
3. /yearly_stats does the same thing but separated into years
4. /town_stats gets HDB records and sorts them into towns, then giving a series of time based records for each town

Query parameters(optional):
- start: The start date for filtering records used in the model. Defaults to `2018-01`.  
  Example: `?start=2020-01`
- end: The end date for filtering records used in the model. Defaults to `2021-01`.  
  Example: `?end=2021-12`
- towns: A list of towns to filter the records. Can be provided multiple times.  
  Example: `?towns=Ang%20Mo%20Kio&towns=Bedok`
- flatType: Filter by flat type (e.g., `4-Room`, `5-Room`).  
  Example: `?flatType=4-Room`

### Predicting
1. /linear_regression gives you a linear regression formula and predictions by a linear regression model
2. /polynomial_regression gives you the same thing but for a polynomial regression model
3. /holt_winters gives you predictions according the Holt-Winters model
4. /town_predictions gives the expected Return of Investment based on different towns after a set number of time

Query Parameter(optional):

- start: The start date for filtering records used in the model. Defaults to `2018-01`.  
  Example: `?start=2020-01`
- end: The end date for filtering records used in the model. Defaults to `2021-01`.  
  Example: `?end=2021-12`
- towns: A list of towns to filter the records. Can be provided multiple times.  
  Example: `?towns=Ang%20Mo%20Kio&towns=Bedok`
- flatType: Filter by flat type (e.g., `4-Room`, `5-Room`).  
  Example: `?flatType=4-Room`
- timeAhead: The number of time units (months or years) to predict into the future. Defaults to `5` if not specified.  
  Example: `?timeAhead=12`
- dateBasis: Defines the time unit for prediction:  
  - `monthly` *(default)*: Predicts monthly values.  
  - `yearly`: Predicts yearly values.  
  Example: `?dateBasis=yearly`


