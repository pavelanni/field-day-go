if id "registration" &>/dev/null; then
    echo 'User registration exists'
else
    echo 'Creating user registration'
    useradd -p nfdfhyMCZWk6w -c "Registration" -m registration
fi
