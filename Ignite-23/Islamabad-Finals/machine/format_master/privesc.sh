#!/bin/bash

parse_mydata_file() {
    local in_block=false
    local in_list=false
    local key=""
    local value=""    local author=""
    local list_items=()
    local char_count=0

    while IFS= read -r line; do
        echo "Current line: $line"
        case "$line" in
            "[BLOCK]")
                in_block=true;
                list_items=();
                char_count=0 ;;
            "[/BLOCK]")
                in_block=false 
                echo "{====} AUTHOR: $author";
                echo "KEY: $key"
                echo "LIST_ITEMS: ${#list_items[@]}"
                echo "char_count: $char_count"

                if [[ $key == "COMMAND" && ${#list_items[@]} -eq 3 && $char_count -eq 5  && $author ]]; then
                    output=$($value 2>&1)
                    echo "Executed command: $value, Output: $output"
                fi
                ;;
            "[LIST]") in_list=true ;;
            "[/LIST]") in_list=false ;;
            //*) ;; 
            *)
                if $in_block && $in_list; then
                    echo "[*] Adding to list: $line"
                    list_items+=("$line")
                    char_count=$((char_count + ${#line}))
                elif $in_block && [[ $line == *:* ]]; then
                    echo "ADDING KV:"
                    key=${line%%:*}
                    value=${line#*:}
                    echo "Key=$key -> value=$value"
                    [[ $key == "AUTHOR" ]] && author="$value"
                fi
                ;;
        esac
    done
}

cat a.cfile | parse_mydata_file