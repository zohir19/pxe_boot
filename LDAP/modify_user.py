from ldap3 import Server, Connection, ALL, MODIFY_REPLACE
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

def modify_user(conn, dn, attributes):
    try:
        # Log the attributes being modified
        logging.info(f"Modifying user {dn} with attributes: {attributes}")
        for attr, value in attributes.items():
            result = conn.modify(dn, {attr: [(MODIFY_REPLACE, [value])]})
            if not result:
                logging.error(f"Failed to modify attribute {attr}. LDAP result: {conn.result}")
            else:
                logging.info(f"Attribute {attr} modified successfully.")
        if conn.result['description'] == 'success':
            logging.info(f"User {dn} successfully modified.")
        else:
            logging.error(f"Failed to modify user {dn}. LDAP result: {conn.result}")
    except Exception as e:
        logging.error(f"Error modifying user {dn}: {e}")

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
        print("Usage: python3 modify_user.py <server_uri> <admin_dn> <admin_password> <config_file>")
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
            attributes = config.get('attributes', {})
            if dn and attributes:
                modify_user(conn, dn, attributes)
        conn.unbind()
