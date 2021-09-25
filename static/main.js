const root = document.getElementById('root');

const configApp = {
    login: {
        href: '/login',
        name: 'Авторизация',
        open: loginPage,
    },
    signup: {
        href: '/signup',
        name: 'Регистрация',
        open: signupPage,
    },
    createProfile: {
        href: '/createProfile',
        name: 'Создать профиль',
        open: createProfilePage,
    },
    profile: {
        href: '/profile',
        name: 'Профиль',
        open: profilePage,
    }
}

const emailRegExp = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
const passwordRegExp = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[a-zA-Z]).{8,}$/;

function ajax(method, url, body = null, callback) {
    const xhr = new XMLHttpRequest();
    xhr.open(method, url, true);
    xhr.withCredentials = true;

    xhr.addEventListener('readystatechange', function() {
        if (xhr.readyState !== XMLHttpRequest.DONE) return;

        callback(xhr.status, xhr.responseText);
    });

    if (body) {
        xhr.setRequestHeader('Content-type', 'application/json; charset=utf8');
        xhr.send(JSON.stringify(body));
        return;
    }

    xhr.send();
}

function createInput(type, text, name) {
    const input = document.createElement('input');
    input.type = type;
    input.name = name;
    input.placeholder = text;
    // input.classList.add('form-field');

    return input;
}

function createCenterContainerRow() {
    const divContainer = document.createElement('div');
    divContainer.classList.add('center-container-row');

    return divContainer;
}

function createCenterContainerColumn() {
    const divContainer = document.createElement('div');
    divContainer.classList.add('center-container-column');

    return divContainer;
}

function loginPage() {
    root.innerHTML = '';

    // --------------------------------------------------------
    window.addEventListener('load', (e) => {
        e.preventDefault();
        const requestOptions = {
            method: 'GET',
            headers: { 'Content-Type': 'application/json'},
        };
        fetch("http://127.0.0.1:8080/api/v1/cookie", requestOptions)
        .then(response => 
                response.json().then(data => ({
                    data: data,
                    status: response.status
                })
            ).then(res => {
                if (res.status === 200 && res.data.status === 200) {
                    profilePage();
                }
                console.log(res.data);
            })).catch((error) => console.log(error));
    })
    // --------------------------------------------------------

    const header = createCenterContainerRow();

    const headerText = document.createElement('span');
    headerText.textContent = 'Войти';
    headerText.classList.add('login-header');

    header.appendChild(headerText);
    root.appendChild(header);

    const form = document.createElement('form');
    form.classList.add('login-form');

    const emailInput = createInput('email', 'Почта', 'email');
    emailInput.addEventListener('input', () => {
        const test = emailInput.value.length === 0 || emailRegExp.test(emailInput.value);
        
        if (test) {
            emailInput.className = 'form-field-valid';
            // emailFieldWithError.removeChild(emailInputError);
        } else {
            emailInput.className = 'form-field-novalid'
            // passwordFieldWithError.removeChild(passwordInputError);
        }
    })

    const passwordInput = createInput('password', 'Пароль', 'password');
    passwordInput.addEventListener('input', () => {
        const test = passwordInput.value.length === 0 || passwordRegExp.test(passwordInput.value);
        
        if (test) {
            passwordInput.className = 'form-field-valid';
        } else {
            passwordInput.className = 'form-field-novalid'
        }
    })

    window.addEventListener('load', () => {
            const testEmail = emailInput.value.length === 0 || emailRegExp.test(emailInput.value);
            emailInput.className = testEmail ? 'form-field-valid' : 'form-field-novalid';
            const testPassword = passwordInput.value.length === 0 || passwordRegExp.test(passwordInput.value);
            passwordInput.className = testPassword ? 'form-field-valid' : 'form-field-novalid';
        })
    
    // кнопка войти
    const submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.classList.add('login-button');

    const buttonFilling = createCenterContainerRow();
    const buttonText = document.createElement('span');
    buttonText.textContent = 'Войти';
    buttonText.classList.add('login-button-text');
    const buttonIcon = document.createElement('img');
    buttonIcon.src = './svg/next.svg';
    buttonIcon.classList.add('svg-next');

    buttonFilling.appendChild(buttonText);
    buttonFilling.appendChild(buttonIcon);

    submitButton.appendChild(buttonFilling);

    const emailIcon = document.createElement('img');
    emailIcon.src = './svg/email.svg'
    emailIcon.classList.add('input-icon');
    const passwordIcon = document.createElement('img');
    passwordIcon.src = './svg/password.svg';
    passwordIcon.classList.add('input-icon');

    const emailFieldWithIcon = document.createElement('div');
    emailFieldWithIcon.classList.add('input-with-icon');
    const passwordFieldWithIcon = document.createElement('div');
    passwordFieldWithIcon.classList.add('input-with-icon');

    const logoBg = document.createElement('div');
    logoBg.classList.add('drip-logo-bg');

    const formContainer = createCenterContainerRow();

    emailFieldWithIcon.appendChild(emailInput);
    emailFieldWithIcon.appendChild(emailIcon);
    passwordFieldWithIcon.appendChild(passwordInput);
    passwordFieldWithIcon.appendChild(passwordIcon);

    // const emailFieldWithError = createCenterContainerColumn();
    // const passwordFieldWithError = createCenterContainerColumn();

    // emailFieldWithError.appendChild(emailFieldWithIcon);
    // passwordFieldWithError.appendChild(passwordFieldWithIcon);

    // logoBg.appendChild(emailFieldWithError);
    // logoBg.appendChild(passwordFieldWithError);

    logoBg.appendChild(emailFieldWithIcon);
    logoBg.appendChild(passwordFieldWithIcon);

    form.appendChild(logoBg);
    form.appendChild(submitButton);

    formContainer.appendChild(form);

    const regLinkContainer = createCenterContainerRow();
    const regLink = document.createElement('a');
    regLink.classList.add('reg-link');
    regLink.href = '/signup';
    regLink.textContent = 'Зарегестрироваться';
    regLink.dataset.section = 'signup';
    regLinkContainer.appendChild(regLink);

    form.addEventListener('submit', (e) => {
        e.preventDefault();
        const testEmail = emailRegExp.test(emailInput.value);
        const testPassword = passwordRegExp.test(passwordInput.value);

        if (!testEmail) {
            emailInput.className = 'form-field-novalid';
            // emailInputError = document.createElement('span')
            // emailInputError.textContent = 'error'
            // emailFieldWithError.appendChild(emailInputError)
        }
        
        if (!testPassword) {
            passwordInput.className = 'form-field-novalid';
            // passwordInputError = document.createElement('span')
            // passwordInputError.textContent = 'error'
            // passwordFieldWithError.appendChild(passwordInputError)
        }
        if (!testEmail || !testPassword) {
            return;
        }

        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(
                { 
                    'email': email, 
                    'password': password,
                }
                )
        };
        fetch("http://127.0.0.1:8080/api/v1/login", requestOptions)
        .then(response => 
                response.json().then(data => ({
                    data: data,
                    status: response.status
                })
            ).then(res => {
                if (res.status === 200 && res.data.status === 200) {
                    profilePage();
                } else if (res.data.status === 404) {
                    const userNotFound = document.createElement('span')
                    userNotFound.textContent = 'Вы не зарегестрированы'
                    form.appendChild(userNotFound)
                }
            })).catch((error) => console.log(error));
    })

    root.appendChild(formContainer);
    root.appendChild(regLinkContainer);
}

function signupPage() {
    root.innerHTML = '';

    const header = createCenterContainerRow();

    const headerText = document.createElement('span');
    headerText.textContent = 'Регистрация';
    headerText.classList.add('login-header');

    header.appendChild(headerText);
    root.appendChild(header);

    const form = document.createElement('form');
    form.classList.add('login-form');

    const emailInput = createInput('email', 'Почта', 'email');
    const passwordInput = createInput('password', 'Пароль', 'password');
    const repeatPasswordInput = createInput('password', 'Повторите пароль', 'password');

    // кнопка зарегестрироваться
    const submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.classList.add('login-button');

    const buttonFilling = createCenterContainerRow();
    const buttonText = document.createElement('span');
    buttonText.textContent = 'Зарегистрироваться';
    buttonText.classList.add('login-button-text');
    const buttonIcon = document.createElement('img');
    buttonIcon.src = './svg/next.svg';
    buttonIcon.classList.add('svg-next');

    buttonFilling.appendChild(buttonText);
    buttonFilling.appendChild(buttonIcon);

    submitButton.appendChild(buttonFilling);

    const emailIcon = document.createElement('img');
    emailIcon.src = './svg/email.svg'
    emailIcon.classList.add('input-icon');
    const passwordIcon = document.createElement('img');
    passwordIcon.src = './svg/password.svg';
    passwordIcon.classList.add('input-icon');
    const repeatPasswordIcon = document.createElement('img');
    repeatPasswordIcon.src = './svg/password.svg';
    repeatPasswordIcon.classList.add('input-icon');

    const emailFieldWithIcon = document.createElement('div');
    emailFieldWithIcon.classList.add('input-with-icon');
    const passwordFieldWithIcon = document.createElement('div');
    passwordFieldWithIcon.classList.add('input-with-icon');
    const repeatPasswordFieldWithIcon = document.createElement('div');
    repeatPasswordFieldWithIcon.classList.add('input-with-icon');

    const logoBg = document.createElement('div');
    logoBg.classList.add('drip-logo-bg');

    const formContainer = createCenterContainerRow();

    emailFieldWithIcon.appendChild(emailInput);
    emailFieldWithIcon.appendChild(emailIcon);
    passwordFieldWithIcon.appendChild(passwordInput);
    passwordFieldWithIcon.appendChild(passwordIcon);
    repeatPasswordFieldWithIcon.appendChild(repeatPasswordInput);
    repeatPasswordFieldWithIcon.appendChild(repeatPasswordIcon);

    logoBg.appendChild(emailFieldWithIcon);
    logoBg.appendChild(passwordFieldWithIcon);
    logoBg.appendChild(repeatPasswordFieldWithIcon);

    form.appendChild(logoBg);
    form.appendChild(submitButton);

    formContainer.appendChild(form);

    const regLink = document.createElement('a');
    regLink.href = '/createProfile';
    regLink.textContent = 'Зарегестрироваться';
    regLink.dataset.section = 'createProfile';

    form.addEventListener('submit', (e) => {
        e.preventDefault();

        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();

        ajax(
            'POST',
            '/signup',
            {email, password},
            (status) => {
                if (status === 200) {
                    profilePage();
                    return;
                }
                alert('Authorized error');
            }
        );
    })

    root.appendChild(formContainer);
}

function createProfilePage() {
    root.innerHTML = '<span>Create Profile Page</span>';
}

function profilePage() {
    root.innerHTML = '<span>Profile Page</span>';
}

loginPage();


root.addEventListener('click', (e) => {
    const {target} = e;

    if (target instanceof HTMLAnchorElement) {
        e.preventDefault();

        configApp[target.dataset.section].open();
    }
})
