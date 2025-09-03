Gator -- A Blog Aggregator

Requirements: 
    Postgres and Go must be installed. 

Installation: 
    Run ``go install gator`` 

Set Up: 
    Create a text file called ".gatorconfig.json" in your home directory.
    It should include a json entry with your postgres url and look something like this:
        {"db_url":"postgres://postgres:9374@localhost:5432/gator?sslmode=disable",}


Commands: 
    register (user): creates a new user and logs that user in

    login (user): logs in as an existing user

    reset:   clears all stored data (users, tables, etc )

    agg: scrapes all current feeds for posts and stores them in the db

    addfeed (url): adds a feed to current user's following list.

    feeds:  prints current user's following list

    unfollow (url): removes a feed from the current user's following list

    browse: prints previously scraped posts for the current user
