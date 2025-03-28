#!/bin/bash

# Keycloak API details
KEYCLOAK_SERVER="http://192.168.56.100:8080"
REALM="openinnovation"
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJtaUtoMUJwNmtCVHpaVnF3clFMb25panJNaGtuNGNmNk9WenNQQzVpVWVzIn0.eyJleHAiOjE3NDMxMzYwNDQsImlhdCI6MTc0MzEwMDA0NCwianRpIjoiMDdkNDM4YmUtMTcyYy00Nzc0LWIwZGItOTg4MzNhZjk5ZjUzIiwiaXNzIjoiaHR0cDovLzE5Mi4xNjguNTYuMTAwOjgwODAvcmVhbG1zL29wZW5pbm5vdmF0aW9uIiwiYXVkIjoicmVhbG0tbWFuYWdlbWVudCIsInN1YiI6ImJjMTFkNTdiLWVkZjUtNDFlNC1iZWRlLTAyM2M5MjFjNjhmMyIsInR5cCI6IkJlYXJlciIsImF6cCI6InNzaC1sb2dpbiIsInNpZCI6ImU4MTkyZjdlLTIzM2YtNGFkZi1iMDFiLTVhZjUyZWM2ZjIyYiIsImFjciI6IjEiLCJyZXNvdXJjZV9hY2Nlc3MiOnsicmVhbG0tbWFuYWdlbWVudCI6eyJyb2xlcyI6WyJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsInZpZXctcmVhbG0iLCJtYW5hZ2UtaWRlbnRpdHktcHJvdmlkZXJzIiwiaW1wZXJzb25hdGlvbiIsInJlYWxtLWFkbWluIiwiY3JlYXRlLWNsaWVudCIsIm1hbmFnZS11c2VycyIsInF1ZXJ5LXJlYWxtcyIsInZpZXctYXV0aG9yaXphdGlvbiIsInF1ZXJ5LWNsaWVudHMiLCJxdWVyeS11c2VycyIsIm1hbmFnZS1ldmVudHMiLCJtYW5hZ2UtcmVhbG0iLCJ2aWV3LWV2ZW50cyIsInZpZXctdXNlcnMiLCJ2aWV3LWNsaWVudHMiLCJtYW5hZ2UtYXV0aG9yaXphdGlvbiIsIm1hbmFnZS1jbGllbnRzIiwicXVlcnktZ3JvdXBzIl19fSwic2NvcGUiOiJwcm9maWxlIGVtYWlsIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJhbmVzIGxvdWRqZWRpIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiYW5lcyIsImdpdmVuX25hbWUiOiJhbmVzIiwiZmFtaWx5X25hbWUiOiJsb3VkamVkaSIsImVtYWlsIjoibG91ZGplZGlhbmVzQGdtYWlsLmNvbSJ9.GYT953D0TLsM11XLiU-WvvPch2TUsQoS9LKE-dJnWHEuQVL2fmYf1aYoc009V5oEXkwmubMYJU6ehz5i4E-o8FjMDm3DRo_BXKtWiVsbEoZYZ0xZVnxGx8ZnWdVfZ4pS-0IRvy-GW_PgEsE3-RlZ2oDr1C4mW3aIOgRxE75ChNLfvfvqnsZCgN6UCB21Gi50fOxZzAc5CL0vkIEr8489Jryi4tRyma9QyJAJA-tQfGq3Ru8Gb0lSjwcC0ULgrJLsIvRzsStUQ7pw1qmZmUJZCRmDglgV84kMwe0WmM_FzRE2M511hzFGGEPcC6UpvTi2gJBFvm6gXp6HgqpgSfZoQQ"

get_keycloak_user_id() {
    local USERNAME=$1
    local USER_ID
    USER_ID=$(curl -s -X GET "$KEYCLOAK_SERVER/admin/realms/$REALM/users?username=$USERNAME&exact=true" \
        -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')
    if [[ "$USER_ID" == "null" || -z "$USER_ID" ]]; then
        return 1  # Return a non-zero exit status to indicate failure
    else
        echo "$USER_ID"  # Output the USER_ID
    fi
}

get_keycloak_user_groups() {
    local USER_ID=$1
    #echo "$USER_ID"
    #echo "getting into functio"
    RESPONSE=$(curl -s -X GET "$KEYCLOAK_SERVER/admin/realms/$REALM/users/$USER_ID/groups" -H "Authorization: Bearer $TOKEN" | jq -r '.[].name')
    echo "$RESPONSE"





    #echo "$RESPONSE"

    #echo "$TESTGROUP"
    #if [[ -z "$GROUPS" ]]; then
    #    echo "No groups found for user ID: $USER_ID"
    #else
    #    echo "Groups for user ID $USER_ID: $GROUPS"
    #fi
}

while true; do
    # Get users from Keycloak
    USERNAMES=$(curl -s -X GET "$KEYCLOAK_SERVER/admin/realms/$REALM/users" \
        -H "Authorization: Bearer $TOKEN" | jq -r '.[].username')

    if [[ -z "$USERNAMES" ]]; then
        echo "Failed to retrieve users or no users found."
        exit 1
    fi

    # Filter out users that already exist in Bright
    EXISTING_USERS=()
    NEW_USERS=()
    for USERNAME in $USERNAMES; do
        if cmsh -c "user list" | grep -w "$USERNAME" &>/dev/null; then
            EXISTING_USERS+=("$USERNAME")
        else
            NEW_USERS+=("$USERNAME")
        fi
    done

    # Main Menu
    echo "Choose an option:"
    select OPTION in "Import new users" "Delete Users" "View Existing Users" "Quit"; do
        case $OPTION in
            "Import new users")
                while true; do
                    if [[ ${#NEW_USERS[@]} -eq 0 ]]; then
                        echo "No new users to import."
                        break
                    fi

                    echo "Available users to import:"
                    select USERNAME in "${NEW_USERS[@]}" "Back to Main Menu"; do
                        if [[ "$USERNAME" == "Back to Main Menu" ]]; then
                            break 2
                        elif [[ -n "$USERNAME" ]]; then
                            echo "You selected: $USERNAME"

                            # Find the next available UID in Bright
                            HIGHEST_UID=$(cmsh -c "user; list" | awk '{print $2}' | sort -n | tail -1)
                            NEXT_UID=$((HIGHEST_UID + 1))

                            echo "The next available UID is: $NEXT_UID"
                            read -p "Do you want to use this UID? (y/n): " CONFIRM_UID

                            if [[ "$CONFIRM_UID" != "y" ]]; then
                                read -p "Enter the UID you want: " CUSTOM_UID
                                UID_TO_USE=$CUSTOM_UID
                            else
                                UID_TO_USE=$NEXT_UID
                            fi

                            read -p "Are you sure you want to create user '$USERNAME' with UID $UID_TO_USE? (y/n): " CONFIRM_CREATE
                            if [[ "$CONFIRM_CREATE" == "y" ]]; then
                                cmsh -c "user; add $USERNAME; set id $UID_TO_USE; commit"
				KUID=$(get_keycloak_user_id "$USERNAME")
				echo "the $USERNAME keykloack ID is $KUID"
                                echo "now getting the groups of $USERNAME"
				USERSGROUP=$(get_keycloak_user_groups "$KUID")
				echo "$USERSGROUP" | while IFS= read -r GROUPP; do
					if [[ -z "$GROUPP" ]]; then
                                                continue  # Skip empty lines
                                        fi
					GROUP_NAME=$(echo "$GROUPP" | tr ' ' '_')

					# Check if group exists
					if ! cmsh -c "group; list" | grep -qw "$GROUP_NAME"; then
						echo "Group '$GROUP_NAME' does not exist. Creating it in Bright..."
						cmsh -c "group; add $GROUP_NAME; commit"
					fi

					# Add user to the group
                                        echo "Adding $USERNAME to group '$GROUP_NAME'"
                                        cmsh -c "group; append $GROUP_NAME members $USERNAME; commit"
				done

                                if [[ $? -eq 0 ]]; then
                                    echo "User $USERNAME created successfully!"
                                    NEW_USERS=("${NEW_USERS[@]/$USERNAME}")
                                else
                                    echo "Failed to create user $USERNAME."
                                fi
                            else
                                echo "Operation canceled."
                            fi
                        else
                            echo "Invalid selection. Please choose a valid user."
                        fi
                    done
                done
                ;;

            "Delete Users")
                while true; do
                    if [[ ${#EXISTING_USERS[@]} -eq 0 ]]; then
                        echo "No existing users found."
                        break
                    fi

                    echo "Existing users in Bright:"
                    select USERNAME in "${EXISTING_USERS[@]}" "Back to Main Menu"; do
                        if [[ "$USERNAME" == "Back to Main Menu" ]]; then
                            break 2
                        elif [[ -n "$USERNAME" ]]; then
                            echo "You selected: $USERNAME"
                            read -p "Do you want to delete this user? (y/n): " CONFIRM_DELETE
                            if [[ "$CONFIRM_DELETE" == "y" ]]; then
                                cmsh -c "user; remove $USERNAME; commit"
				if [[ $? -eq 0 ]]; then
                                    echo "User $USERNAME deleted successfully!"
                                    EXISTING_USERS=("${EXISTING_USERS[@]/$USERNAME}")
                                else
                                    echo "Failed to delete user $USERNAME."
                                fi
				sudo chown nobody:nobody /home/$USERNAME
				echo "Changed the deleted user : $USERNAME home directory ownership to nobody:nobdy  to avoid confusion and keep the data of this useri"
				sudo sss_cache -u $USERNAME
				sudo sss_cache -g $USERNAME
				echo "Cleared the user and group sss cache for $USERNAME"
                            else
                                echo "Operation canceled."
                            fi
                        else
                            echo "Invalid selection. Please choose a valid user."
                        fi
                    done
                done
                ;;
            "View Existing Users")
                echo "The Existing new users on the cluster:"
                cmsh -c "user; list; quit;"

                while true; do
                    read -p "Do you want to go back to the main menu? (y/n): " response
                    case $response in
                        [Yy]*)
                            break 2 # Exit from the inner loop and return to main menu
                            ;;
                        [Nn]*)
                            echo "Exiting."
                            exit 0
                            ;;
                        *)
                            echo "Please answer yes or no."
                            ;;
                    esac
                done
                ;;

            "Quit")
                echo "Exiting."
                exit 0
                ;;

            *)
                echo "Invalid option. Please select again."
                ;;
        esac
        break
    done

done

