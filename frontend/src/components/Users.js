import React, { useState, useEffect } from "react";
import axios from "axios";
import {useNavigate} from "react-router-dom";
import constants from "../constants";

const Users = () => {
    const { API_URL } = constants;
    const navigate = useNavigate(); // Hook to handle navigation
    const [users, setUsers] = useState([]);
    const [showAddUserModal, setShowAddUserModal] = useState(false);
    const [newUser, setNewUser] = useState({
        name: "",
        email: "",
        password: "",
        nasClientIP: "",
        role: "ROLE_USER",
    });

    const fetchUsers = async () => {
        try {
            const authToken = localStorage.getItem("auth_token");
            if (!authToken) {
                navigate('/login'); // Redirect to login if no token is found
                return;
            }
            const response = await axios.get(`${API_URL}/api/v1/users`, {
                headers: {
                    Authorization: `${authToken}`,
                },
            });
            setUsers(response.data.data); // Assuming the response contains a 'data' array
        } catch (error) {
            console.error("Error fetching users:", error);
        }
    };

    // Fetch users on page load
    useEffect(() => {
        fetchUsers();
    }, []);

    // Delete user
    const handleDelete = async (userId) => {
        const confirmDelete = window.confirm("Are you sure you want to delete this user?");
        if (!confirmDelete) return;

        try {
            const authToken = localStorage.getItem("auth_token");
            await axios.delete(`${API_URL}/api/v1/users/${userId}`, {
                headers: {
                    Authorization: `${authToken}`,
                },
            });

            // Update state to remove the deleted user
            setUsers(users.filter((user) => user.id !== userId));
            alert("User deleted successfully.");
        } catch (error) {
            console.error("Error deleting user:", error);
            alert("Failed to delete user.");
        }
    };

    // Handle form input change
    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setNewUser((prevUser) => ({
            ...prevUser,
            [name]: value,
        }));
    };

    // Handle add user submit
    const handleAddUser = async () => {
        try {
            const authToken = localStorage.getItem("auth_token");
            const response = await axios.post(
                `${API_URL}/api/v1/users`,
                {
                    name: newUser.name,
                    email: newUser.email,
                    password: newUser.password,
                    nasClientIP: newUser.nasClientIP,
                    role: newUser.role,
                },
                {
                    headers: {
                        Authorization: `${authToken}`,
                    },
                }
            );

            fetchUsers();

            alert("User added successfully.");
            setShowAddUserModal(false); // Close modal after adding user
        } catch (error) {
            console.error("Error adding user:", error);
            alert("Failed to add user.");
        }
    };

    return (
        <div style={{ padding: "20px" }}>
            <h1 style={{ marginBottom: "20px" }}>Users</h1>

            {/* Add User Button */}
            <button
                onClick={() => setShowAddUserModal(true)}
                style={{
                    marginBottom: "20px",
                    padding: "10px 20px",
                    backgroundColor: "#28a745",
                    color: "#fff",
                    border: "none",
                    borderRadius: "5px",
                    cursor: "pointer",
                }}
            >
                Add User
            </button>

            {/* User Table */}
            <table
                style={{
                    width: "100%",
                    borderCollapse: "collapse",
                    marginBottom: "20px",
                }}
            >
                <thead>
                <tr
                    style={{
                        backgroundColor: "#f4f4f4",
                        borderBottom: "2px solid #ddd",
                        textAlign: "left",
                    }}
                >
                    <th style={{ padding: "12px" }}>Name</th>
                    <th style={{ padding: "12px" }}>Email</th>
                    <th style={{ padding: "12px" }}>NAS Client IP</th>
                    <th style={{ padding: "12px" }}>Role</th>
                    <th style={{ padding: "12px", textAlign: "center" }}>Actions</th>
                </tr>
                </thead>
                <tbody>
                {users.map((user) => (
                    <tr
                        key={user.id}
                        style={{
                            borderBottom: "1px solid #ddd",
                        }}
                    >
                        <td style={{ padding: "12px" }}>{user.name}</td>
                        <td style={{ padding: "12px" }}>{user.email}</td>
                        <td style={{ padding: "12px" }}>{user.nasClientIP}</td>
                        <td style={{ padding: "12px" }}>
                            {user.role === "ROLE_ADMIN" ? "Admin" : "User"}
                        </td>
                        <td
                            style={{
                                padding: "12px",
                                textAlign: "center",
                            }}
                        >
                            <button
                                onClick={() => handleDelete(user.id)}
                                style={{
                                    backgroundColor: "#dc3545",
                                    color: "#fff",
                                    border: "none",
                                    borderRadius: "5px",
                                    cursor: "pointer",
                                    padding: "5px 10px",
                                }}
                            >
                                Delete
                            </button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>

            {/* Add User Modal */}
            {showAddUserModal && (
                <div
                    style={{
                        position: "fixed",
                        top: "50%",
                        left: "50%",
                        transform: "translate(-50%, -50%)",
                        backgroundColor: "#fff",
                        padding: "30px",
                        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.2)",
                        borderRadius: "8px",
                        width: "300px",
                        zIndex: "1000",
                    }}
                >
                    <h2>Add User</h2>
                    <label>Name:</label>
                    <input
                        type="text"
                        name="name"
                        value={newUser.name}
                        onChange={handleInputChange}
                        style={{ marginBottom: "10px", width: "100%", padding: "8px" }}
                    />
                    <br />
                    <label>Email:</label>
                    <input
                        type="email"
                        name="email"
                        value={newUser.email}
                        onChange={handleInputChange}
                        style={{ marginBottom: "10px", width: "100%", padding: "8px" }}
                    />
                    <br />
                    <label>Password:</label>
                    <input
                        type="password"
                        name="password"
                        value={newUser.password}
                        onChange={handleInputChange}
                        style={{ marginBottom: "10px", width: "100%", padding: "8px" }}
                    />
                    <br />
                    <label>NAS Client IP:</label>
                    <input
                        type="text"
                        name="nasClientIP"
                        value={newUser.nasClientIP}
                        onChange={handleInputChange}
                        style={{ marginBottom: "10px", width: "100%", padding: "8px" }}
                    />
                    <br />
                    <label>Role:</label>
                    <select
                        name="role"
                        value={newUser.role}
                        onChange={handleInputChange}
                        style={{ marginBottom: "10px", width: "100%", padding: "8px" }}
                    >
                        <option value="ROLE_USER">User</option>
                        <option value="ROLE_ADMIN">Admin</option>
                    </select>
                    <br />
                    <button
                        onClick={handleAddUser}
                        style={{
                            marginTop: "20px",
                            padding: "10px 20px",
                            backgroundColor: "#28a745",
                            color: "#fff",
                            border: "none",
                            borderRadius: "5px",
                            cursor: "pointer",
                        }}
                    >
                        Add User
                    </button>
                    <button
                        onClick={() => setShowAddUserModal(false)}
                        style={{
                            marginTop: "20px",
                            padding: "10px 20px",
                            backgroundColor: "#dc3545",
                            color: "#fff",
                            border: "none",
                            borderRadius: "5px",
                            cursor: "pointer",
                        }}
                    >
                        Cancel
                    </button>
                </div>
            )}
        </div>
    );
};

export default Users;
