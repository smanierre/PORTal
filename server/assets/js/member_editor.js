document.addEventListener("alpine:init", () => {
    Alpine.data("memberEditor", () => ({
        password: "",
        confirmPassword: "",
        changes: false,
        passwordChanges: false,
        validPassword: false,
        submitDisabled: true,
    }))
})