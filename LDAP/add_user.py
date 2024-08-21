from ldap3 import Server, Connection, ALL
import json
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

def add_user(conn, dn, attributes):
    try:
        conn.add(dn, attributes=attributes)
        if conn.result['description'] == 'success':
            logging.info(f"User {dn} successfully added.")
        else:
            logging.error(f"Failed to add user {dn}. LDAP result: {conn.result}")
    except Exception as e:
        logging.error(f"Error adding user {dn}: {e}")

def load_config(file_path):
    try:
        with open(file_path, 'r') as file:
            config = json.load(file)
        logging.info("Configuration loaded successfully.")
        return config
    except Exception as e:
        logging.error(f"Error loading configuration file: {e}")
        return None

if __name__ == "__main__":
    if len(sys.argv) != 5:
        print("Usage: python3 add_user.py <server_uri> <admin_dn> <admin_password> <config_file>")
        sys.exit(1)

    server_uri = sys.argv[1]
    admin_dn = sys.argv[2]
    admin_password = sys.argv[3]
    config_file = sys.argv[4]

    conn = connect_to_ldap(server_uri, admin_dn, admin_password)
    if conn:
        config = load_config(config_file)
        if config:
            dn = config.get('dn')
            attributes = {
                'objectClass': ['inetOrgPerson', 'posixAccount'],
                'cn': config.get('cn'),
                'sn': config.get('sn'),
                'uid': config.get('uid'),
                'uidNumber': config.get('uidNumber'),
                'gidNumber': config.get('gidNumber'),
                'homeDirectory': config.get('homeDirectory'),
                'loginShell': config.get('loginShell'),
                'userPassword': config.get('userPassword')
            }
            if dn:
                add_user(conn, dn, attributes)
        conn.unbind()
