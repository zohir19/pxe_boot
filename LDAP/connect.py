# import class and constants
from ldap3 import Server, Connection, ALL

# define the server
s = Server('ldap://192.168.56.6', get_info=ALL)  # define an unsecure LDAP server, requesting info on DSE and schema

# define the connection
c = Connection(s, user='cn=admin,dc=hpcme,dc=com', password='init12')

# perform the Bind operation
if not c.bind():
    print('error in bind', c.result)
~
