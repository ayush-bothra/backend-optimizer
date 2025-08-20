package service


/*
this is the business logic layer
possibly the more important one of
all the files present in this folder

will integrate the logic layer with 
the cache and DB and the AI model

imports needed:
will import the cache, db, ai and 
utils packages

functions:
checkPortfolioRisk(portfolioID int) (float64, error)

check in the cache, if miss then check for DB
pass data to AI model from the DB, get the result
result is in float64, store this in the cache, DB
and then output the result
*/