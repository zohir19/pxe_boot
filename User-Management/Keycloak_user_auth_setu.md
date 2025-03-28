
# Setting Up a Centralized User Authentication and Management System with Keycloak

This document outlines the steps to configure a centralized user authentication and management system using Keycloak.

## Step 1: Set Up the Keycloak Server

Keycloak can be installed on the head node or on a separate server, depending on your infrastructure requirements.

### Download Keycloak

Obtain the Keycloak distribution from the official Keycloak website. For this example, we will use Keycloak version 26.1.4.
Download and extract the `keycloak-26.1.4.zip` archive. After extraction, you will have a directory named `keycloak-26.1.4`.

### Start Keycloak

To start the Keycloak server, navigate to the extracted `keycloak-26.1.4` directory from your terminal.

On Linux:
```bash
bin/kc.sh start-dev
```

On Windows:
```bash
bin\kc.bat start-dev
```

> **Note:** Using the `start-dev` option starts Keycloak in development mode. This mode is ideal for testing and initial setup, offering default settings that are convenient for developers. However, it is not recommended for production environments, as it is optimized for ease of use rather than security or performance.

### Create an admin user

Keycloak has no default admin user. You need to create an admin user before you can start Keycloak.

1. Open `http://localhost:8080/`.
2. Fill in the form with your preferred username and password.

### Log in to the Admin Console

1. Go to the Keycloak Admin Console.
2. Log in with the username and password you created earlier.

### Create a realm

A realm in Keycloak is equivalent to a tenant. Each realm allows an administrator to create isolated groups of applications and users. Initially, Keycloak includes a single realm, called `master`. Use this realm only for managing Keycloak and not for managing any applications.

To create your first realm:

1. Open the Keycloak Admin Console.
2. Click `Keycloak` next to the master realm, then click `Create Realm`.
3. Enter `myrealm` in the Realm name field.
4. Click `Create`.

### Create an SSH client for the realm

1. Log in to the Keycloak Administration Console.
2. Select the realm for which you want to create the client.
3. Click on "Clients" from the left-hand menu, and then click on the "Create" button.
4. In the "Client ID" field, enter "ssh-login".
5. Set the "Client Protocol" to "openid-connect".
6. In the "Valid Redirect URIs" field, enter `urn:ietf:wg:oauth:2.0:oob`.
7. In the "Access Type" field, select "confidential".
8. In the "Standard Flow Enabled" field, select "ON".
9. In the "Direct Access Grants Enabled" field, select "ON".
10. Click on the "Save" button to create the client.

#### Verify client settings

1. Go to the client `ssh-login` and confirm the following are set and correct:
    - Client ID: `ssh-login`
    - Valid redirect URIs: `urn:ietf:wg:oauth:2.0:oob`
    - Client authentication: `On`
    - Authorization: `On`
    - Authentication flow: `Standard flow` --> `On`, `Direct access grants` --> `On`
    - Front channel logout: `On`
    - Front-channel logout session required: `On`

### Get client credentials

1. Go to the "Clients" page in the Keycloak Administration Console.
2. Select the "ssh-login" client from the list.
3. Click on the "Credentials" tab.
4. The client secret will be displayed under the "Client Secret" section.

### Import the external AD LDAP (if exists)

1. On the left panel of Keycloak, go to **User federation**.
2. Click on **Add new provider** and choose **LDAP**.
3. In the UI, configure the following:
    - Vendor: Active Directory
    - Connection URL: `ldap://yourserver:389`
    - Enable StartTLS: off
    - Use Truststore SPI: Always
    - Connection pooling: On

4. Test connection. It should say: `Successfully connected to LDAP`.

### Bind settings

- Bind type: simple
- Bind DN: `CN=Administrator,CN=Users,DC=ad,DC=com`
- Bind credentials: "Your_AD_password"

Test Authentication. It should say: `Successfully connected to LDAP`.

### LDAP Searching and Updating

- Edit mode: READ_ONLY
- User DN: `CN=Users,DC=ad,DC=com`
- Username LDAP attribute: `sAMAccountName`
- RDN LDAP attribute: `cn`
- UUID LDAP attribute: `objectGUID`
- User object classes: person, organizationalPerson, user
- Search scope: One Level
- Pagination: On

### Synchronization settings

- Import users: On
- Sync Registration: On

### Kerberos Integration

- Allow Kerberos authentication: Off
- Use Kerberos for password authentication: On

### Enable Group Import and Sync Users into Groups

1. On the **Mappers** tab:
    - Create a new mapper named **group mapper** with the following settings:
      - Name: group mapper
      - Mapper type: group-ldap-mapper
      - LDAP Groups DN: `CN=Users,DC=ad,DC=com`
      - Group Name LDAP Attribute: `cn`
      - Group Object Classes: group
      - Preserve Group Inheritance: Off
      - Ignore Missing Groups: Off
      - Membership LDAP Attribute: `member`
      - Membership Attribute Type: DN
      - Membership User LDAP Attribute: `sAMAccountName`
      - Mode: LDAP_ONLY
      - User Groups Retrieve Strategy: LOAD_GROUPS_BY_MEMBER_ATTRIBUTE
      - Member-Of LDAP Attribute: `memberoOf`

2. Click **Save**.
3. Click on **Action** and then select **Sync LDAP groups to Keycloak**.

Now your Keycloak server is well configured and ready.

## Step 2: Setup Keycloak Authentication (Client Configuration)

To download the `kc-ssh-pam` package, visit the GitHub link: [kc-ssh-pam GitHub](https://github.com/kha7iq/kc-ssh-pam). Navigate to the releases section and select the appropriate package for your operating system architecture.

The `kc-ssh-pam` module is designed to simplify user authentication, allowing seamless access to Linux systems through SSH. This program integrates with Keycloak to acquire a password grant token based on the user's login credentials, which include their username and password. Additionally, if two-factor authentication is enabled for the user, the program also supports the use of an OTP code.

### Installation of the `kc-ssh-pam` Package

#### For DEB

To install the package, execute the following command:
```bash
sudo dpkg -i kc-ssh-pam_amd64.deb
```

#### For RPM

To install the package, execute the following command:
```bash
sudo rpm -i kc-ssh-pam_amd64.rpm
```

### Configuration

Ensure that **SELinux** is set to **Permissive** mode.

To configure the `kc-ssh-pam` configuration file, follow these steps:
```bash
vim /opt/kc-ssh-pam/
```
Add the following configuration:
```bash
realm = "ssh-demo"  
endpoint = "https://keycloak.server.com"  
clientid = "keycloak-client-id"  
clientsecret = "MIKEcHObWmI3V3pF1hcSqC9KEILfLN"  
clientscope = "openid"
```

Edit the `/etc/pam.d/sshd` file and add the following lines at the top:
```bash
auth sufficient pam_exec.so expose_authtok log=/var/log/kc-ssh-pam.log /opt/kc-ssh-pam/kc-ssh-pam
```

Finally, restart `sshd`:
```bash
systemctl restart sshd
```

## Step 3: Import/Delete/Setup Users on Our Cluster Script

In order for the procedure to work, the user needs to be installed locally on your client cluster but keeps the password through Keycloak.

I have set up a script that does the following:

### Keycloak API Integration
- The script interacts with the Keycloak server to manage user data.
- It retrieves a user's unique ID (`get_keycloak_user_id`) based on their username.
- It also fetches the groups associated with a user (`get_keycloak_user_groups`).

### User Import and Management
- The script periodically checks for users in Keycloak and compares them against existing users in Bright.
- New users in Keycloak (not present in Bright) are imported to Bright and given a unique UID.
- After importing users, the script fetches the groups the user belongs to in Keycloak and ensures these groups exist in Bright. If the groups don't exist, it creates them and adds the user to those groups.

### Main Menu
The script displays a menu allowing the user to:
- Import new users: Import users from Keycloak that are not present in Bright.
- Delete Users: Delete users from Bright.
- View Existing Users: View a list of existing users in Bright.
- Quit: Exit the script.

### User Deletion
If a user is deleted from Bright, their home directory ownership is changed to `nobody:nobody` to avoid confusion, and the script clears the cache for the user.

### Detailed Workflow

#### Import Users
- The script finds new users from Keycloak, then asks the user if they want to import them with a new UID.
- It also creates any missing groups in Bright and adds the user to those groups.

#### Delete Users
- The script allows for the deletion of existing users from Bright, ensuring to clear ownership and cache as necessary.

#### View Users
- The script lists users already present in Bright for review.

### Error Handling
- The script checks if users already exist in Bright and ensures the right groups are created before adding users to those groups.
- If any operation fails, appropriate error messages are displayed.

In summary, the script automates the process of synchronizing users between Keycloak and Bright, including adding groups and managing user deletions.

The script is called import_users.sh and it is in the same directory as this documentation 
