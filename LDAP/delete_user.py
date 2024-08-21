from ldap3 import Server, Connection, ALL, SUBTREE
import sys
import logging

logging.basicConfig(level=logging.INFO)

def connect_to_ldap(server_uri, admin_dn, admin_password):
    try:
        server = Server(server_uri, get_info=ALL)
        conn = Connection(server, admin_dn, admin_password, auto_bind=True)
        logging.info("Successfully connected to the LDAP server.")
        return conn
    except Exception as e:
        logging.error(f"Error connecting to LDAP server: {e}")
        return None

def search_user(conn, base_dn, uid):
    try:
        search_filter = f"(uid={uid})"
        conn.search(search_base=base_dn, search_filter=search_filter, search_scope=SUBTREE, attributes=['entryDN'])

        if conn.entries:
            user_dn = conn.entries[0].entry_dn
            logging.info(f"User {uid} found. DN: {user_dn}")
            return user_dn
        else:
            logging.info(f"User {uid} not found.")
            return None
    except Exception as e:
        logging.error(f"Error searching for user {uid}: {e}")
        return None

def delete_user(conn, user_dn):
    try:
        conn.delete(user_dn)
        if conn.result['description'] == 'success':
            logging.info(f"User {user_dn} successfully deleted.")
        else:
            logging.error(f"Failed to delete user {user_dn}. LDAP result: {conn.result}")
    except Exception as e:
        logging.error(f"Error deleting user {user_dn}: {e}")

if __name__ == "__main__":
    if len(sys.argv) != 6:
        print("Usage: python3 search_user.py <server_uri> <admin_dn> <admin_password> <base_dn> <uid>")
        sys.exit(1)

    server_uri = sys.argv[1]
    admin_dn = sys.argv[2]
    admin_password = sys.argv[3]
    base_dn = sys.argv[4]
    uid = sys.argv[5]

    conn = connect_to_ldap(server_uri, admin_dn, admin_password)
    if conn:
        user_dn = search_user(conn, base_dn, uid)
        if user_dn:
            delete_user(conn, user_dn)
        conn.unbind()
