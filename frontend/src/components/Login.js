import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate  } from 'react-router-dom';
import './Login.css';
import constants from "../constants"; // CSS file for styling

const Login = () => {
    const { API_URL } = constants;
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate ();

    const handleLogin = async (e) => {
        e.preventDefault();

        try {
            const response = await axios.post(`${API_URL}/api/v1/auth/login`, { username, password });
            if (response.status === 200) {
                localStorage.setItem('auth_token', response.data.token);
                navigate('/dashboard');
            }
        } catch (err) {
            setError('Login failed. Please check your credentials.');
        }
    };

    return (
        <div className="login-container">
            <h1>easyNAS</h1>
            <h2>Login to Your Account</h2>
            <form onSubmit={handleLogin} className="login-form">
                <div className="input-group">
                    <label>Username</label>
                    <input
                        type="username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                    />
                </div>
                <div className="input-group">
                    <label>Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                {error && <p className="error-message">{error}</p>}
                <button type="submit" className="login-button">Login</button>
            </form>
        </div>
    );
};

export default Login;
