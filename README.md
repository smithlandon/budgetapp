# budgetapp

I forgot to add the DB tables -> will just pseudocode it real quick


User()\
pk ksuid varchar\
// room for lots more...


Account() (1:1 with User)\
pk ksuid varchar\
fk user_id varchar\
lastUdpatedAt timestamp\
accountId varchar

Budget() (1:Many with User)\
pk ksuid varchar\
fk user_id\
category enum varchar\
target int\
actualSpending int


