import React, { useState, useEffect } from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import axios from 'axios';
import constants from "../constants";

const DatasetDetails = () => {
    const { id } = useParams(); // Get dataset ID from URL
    const { API_URL } = constants;
    const [datasetDetails, setDatasetDetails] = useState(null);
    const [permissions, setPermissions] = useState([]); // Permissions data
    const [isToggling, setIsToggling] = useState(false); // Track toggle action
    const [isModalOpen, setIsModalOpen] = useState(false); // Modal state
    const [users, setUsers] = useState([]); // Users for dropdown
    const [selectedUserId, setSelectedUserId] = useState(""); // Selected user ID
    const [selectedPermission, setSelectedPermission] = useState("rw"); // Selected permission
    const [isSubmitting, setIsSubmitting] = useState(false); // Submission state
    const navigate = useNavigate(); // Hook to handle navigation

    const fetchDatasetDetails = async () => {
        const authToken = localStorage.getItem('auth_token'); // Load auth token from localStorage
        if (!authToken) {
            navigate('/login'); // Redirect to login if no token is found
            return;
        }
        try {
            const response = await axios.get(`${API_URL}/api/v1/nas/pools/naspool/datasets/${id}`, {
                headers: {
                    Authorization: `${authToken}`, // Pass token in Authorization header
                },
            });
            if (response.data.status === 'success') {
                setDatasetDetails(response.data.data);
            }
        } catch (error) {
            console.error('Error fetching dataset details:', error);
        }
    };

    const fetchPermissions = async () => {
        const authToken = localStorage.getItem("auth_token");
        if (!authToken) {
            navigate('/login'); // Redirect to login if no token is found
            return;
        }

        try {
            const response = await axios.get(`${API_URL}/api/v1/nas/pools/naspool/datasets/${id}/nfs-share/permissions`, {
                headers: {
                    Authorization: `${authToken}`,
                },
            });
            if (response.data.status === "success") {
                setPermissions(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching permissions:", error);
        }
    };

    const fetchUsers = async () => {
        const authToken = localStorage.getItem('auth_token'); // Load auth token from localStorage
        if (!authToken) {
            navigate('/login'); // Redirect to login if no token is found
            return;
        }
        try {
            const response = await axios.get(`${API_URL}/api/v1/users`, {
                headers: { Authorization: `${authToken}` },
            });
            if (response.data.status === "success") {
                setUsers(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching users:", error);
        }
    };

    useEffect(() => {
        fetchDatasetDetails();
        fetchPermissions();
        fetchUsers();
    }, [id]);

    const toggleShareStatus = async () => {
        if (!datasetDetails) return;

        const authToken = localStorage.getItem('auth_token'); // Load auth token
        const apiUrl = datasetDetails.shareEnabled
            ? `${API_URL}/api/v1/nas/pools/naspool/datasets/${id}/nfs-share` // DELETE method for disabling share
            : `${API_URL}/api/v1/nas/pools/naspool/datasets/${id}/nfs-share`; // POST method for enabling share

        setIsToggling(true);
        try {
            const method = datasetDetails.shareEnabled ? 'delete' : 'post'; // Dynamically determine method
            const response = await axios({
                method,
                url: apiUrl,
                headers: {
                    Authorization: `${authToken}`, // Pass token in Authorization header
                },
            });

            if (response.data.status === 'success') {
                setDatasetDetails({
                    ...datasetDetails,
                    shareEnabled: !datasetDetails.shareEnabled,
                });
                fetchDatasetDetails()
            } else {
                console.error('Error toggling share status:', response.data.message);
            }
        } catch (error) {
            console.error('Error toggling share status:', error);
        } finally {
            setIsToggling(false);
        }
    };

    const handleDelete = async () => {
        // Show confirmation dialog
        const confirmDelete = window.confirm(`Are you sure you want to delete this dataset?`);

        if (!confirmDelete) {
            return; // Exit if the user cancels
        }

        try {
            const authToken = localStorage.getItem("auth_token");
            await axios.delete(`${API_URL}/api/v1/nas/pools/naspool/datasets/${id}`, {
                headers: { Authorization: `${authToken}` },
            });
            navigate(`/dashboard`);
        } catch (error) {
            console.error("Error deleting dataset.", error);
        }
    };

    const handleViewFilesystem = () => {
        navigate(`/dataset/${id}/filesystem`);
    };

    const handleViewSnapshots = () => {
        navigate(`/dataset/${id}/snapshots`);
    };

    const handleAddPermission = async (e) => {
        e.preventDefault();
        const authToken = localStorage.getItem("auth_token");

        const postData = {
            userId: selectedUserId,
            permission: selectedPermission,
        };

        setIsSubmitting(true);
        try {
            const response = await axios.post(`${API_URL}/api/v1/nas/pools/naspool/datasets/${id}/nfs-share/permissions`, postData, {
                headers: { Authorization: `${authToken}` },
            });
            if (response.data.status === "success") {
                fetchPermissions();
                setIsModalOpen(false); // Close modal
            } else {
                console.error("Error adding permission:", response.data.message);
            }
        } catch (error) {
            console.error("Error adding permission:", error);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleDeletePermission = async (permissionId) => {
        const authToken = localStorage.getItem("auth_token");

        if (window.confirm("Are you sure you want to delete this permission?")) {
            try {
                const response = await axios.delete(`${API_URL}/api/v1/nas/pools/naspool/datasets/${id}/nfs-share/permissions/${permissionId}`, {
                    headers: { Authorization: `${authToken}` },
                });
                if (response.data.status === "success") {
                    setPermissions((prev) => prev.filter((perm) => perm.id !== permissionId));
                } else {
                    console.error("Error deleting permission:", response.data.message);
                }
            } catch (error) {
                console.error("Error deleting permission:", error);
            }
        }
    };

    return (
        <div className="dataset-details">
            <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <h2>Dataset Details</h2>
                <button
                    onClick={handleDelete}
                    style={{
                        padding: "10px 20px",
                        backgroundColor: "#f44336",
                        color: "#fff",
                        border: "none",
                        borderRadius: "5px",
                        cursor: "pointer",
                    }}
                >
                    Delete Dataset
                </button>
            </header>
            {datasetDetails ? (
                <div>
                    <p><strong>Name:</strong> {datasetDetails.name}</p>
                    <p><strong>Quota:</strong> {datasetDetails.quota}</p>
                    <p><strong>Used:</strong> {datasetDetails.used}</p>
                    <p><strong>Available:</strong> {datasetDetails.available}</p>
                    <p>
                        <strong>Share Enabled:</strong> {datasetDetails.shareEnabled ? "Yes" : "No"}
                    </p>
                    <button
                        onClick={toggleShareStatus}
                        disabled={isToggling}
                        style={{
                            padding: "10px 20px",
                            backgroundColor: datasetDetails?.shareEnabled ? "#f44336" : "#4caf50",
                            color: "#fff",
                            border: "none",
                            borderRadius: "5px",
                            marginTop: "10px",
                            cursor: "pointer",
                        }}
                    >
                        {isToggling
                            ? "Processing..."
                            : datasetDetails?.shareEnabled
                                ? "Disable Share"
                                : "Enable Share"}
                    </button>
                    <button
                        onClick={handleViewFilesystem}
                        style={{
                            padding: "10px 20px",
                            backgroundColor: "#2196f3",
                            color: "#fff",
                            border: "none",
                            borderRadius: "5px",
                            marginLeft: "10px",
                            cursor: "pointer",
                        }}
                    >
                        View Filesystem
                    </button>
                    <button
                        onClick={handleViewSnapshots}
                        style={{
                            padding: "10px 20px",
                            backgroundColor: "#41a693",
                            color: "#fff",
                            border: "none",
                            borderRadius: "5px",
                            marginLeft: "10px",
                            cursor: "pointer",
                        }}
                    >
                        View Snapshots
                    </button>
                    <button
                        onClick={() => setIsModalOpen(true)}
                        style={{
                            padding: "10px 20px",
                            backgroundColor: "#0146f3",
                            color: "#fff",
                            border: "none",
                            borderRadius: "5px",
                            marginLeft: "10px",
                            cursor: "pointer",
                        }}
                    >
                        Add Permission
                    </button>
                </div>
            ) : (
                <p>Loading dataset details...</p>
            )}

            <h3 style={{ marginTop: "30px" }}>Permissions</h3>
            {permissions.length > 0 ? (
                <table style={{ width: "100%", borderCollapse: "collapse", marginTop: "10px" }}>
                    <thead>
                    <tr>
                        <th style={{ border: "1px solid #ddd", padding: "8px" }}>User</th>
                        <th style={{ border: "1px solid #ddd", padding: "8px" }}>Email</th>
                        <th style={{ border: "1px solid #ddd", padding: "8px" }}>Client IP</th>
                        <th style={{ border: "1px solid #ddd", padding: "8px" }}>Role</th>
                        <th style={{ border: "1px solid #ddd", padding: "8px" }}>Permission</th>
                        <th style={{ border: "1px solid #ddd", padding: "8px" }}>Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {permissions.map((perm) => (
                        <tr key={perm.id}>
                            <td style={{ border: "1px solid #ddd", padding: "8px" }}>{perm.user.name}</td>
                            <td style={{ border: "1px solid #ddd", padding: "8px" }}>{perm.user.email}</td>
                            <td style={{ border: "1px solid #ddd", padding: "8px" }}>{perm.user.nasClientIP}</td>
                            <td style={{ border: "1px solid #ddd", padding: "8px" }}>{perm.user.role === "ROLE_ADMIN" ? "Admin" : "User"}</td>
                            <td style={{ border: "1px solid #ddd", padding: "8px" }}>{perm.permission === "rw" ? "Read & Write" : "Read-Only"}</td>
                            <td style={{ border: "1px solid #ddd", padding: "8px" }}>
                                <button
                                    onClick={() => handleDeletePermission(perm.id)}
                                    style={{
                                        padding: "5px 10px",
                                        backgroundColor: "#f44336",
                                        color: "#fff",
                                        border: "none",
                                        borderRadius: "5px",
                                        cursor: "pointer",
                                    }}
                                >
                                    Delete
                                </button>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            ) : (
                <p>No permissions available.</p>
            )}

            {isModalOpen && (
                <div className="modal" style={{ position: "fixed", top: "0", left: "0", width: "100%", height: "100%", backgroundColor: "rgba(0, 0, 0, 0.5)", display: "flex", justifyContent: "center", alignItems: "center" }}>
                    <div style={{ backgroundColor: "#fff", padding: "20px", borderRadius: "5px", width: "400px" }}>
                        <h3>Add Permission</h3>
                        <form onSubmit={handleAddPermission}>
                            <div style={{ marginBottom: "10px" }}>
                                <label htmlFor="user-select">User</label>
                                <select
                                    id="user-select"
                                    value={selectedUserId}
                                    onChange={(e) => setSelectedUserId(parseInt(e.target.value, 10))}
                                    required
                                    style={{ width: "100%", padding: "8px", marginTop: "5px" }}
                                >
                                    <option value="">Select a user</option>
                                    {users.map((user) => (
                                        <option key={user.id} value={user.id}>
                                            {user.name} ({user.email})
                                        </option>
                                    ))}
                                </select>
                            </div>
                            <div style={{ marginBottom: "10px" }}>
                                <label htmlFor="permission-select">Permission</label>
                                <select
                                    id="permission-select"
                                    value={selectedPermission}
                                    onChange={(e) => setSelectedPermission(e.target.value)}
                                    required
                                    style={{ width: "100%", padding: "8px", marginTop: "5px" }}
                                >
                                    <option value="r">Read-Only</option>
                                    <option value="rw">Read & Write</option>
                                </select>
                            </div>
                            <div style={{ display: "flex", justifyContent: "space-between", marginTop: "20px" }}>
                                <button
                                    type="button"
                                    onClick={() => setIsModalOpen(false)}
                                    style={{
                                        padding: "10px 20px",
                                        backgroundColor: "#f44336",
                                        color: "#fff",
                                        border: "none",
                                        borderRadius: "5px",
                                        cursor: "pointer",
                                    }}
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    style={{
                                        padding: "10px 20px",
                                        backgroundColor: "#4caf50",
                                        color: "#fff",
                                        border: "none",
                                        borderRadius: "5px",
                                        cursor: isSubmitting ? "not-allowed" : "pointer",
                                    }}
                                >
                                    {isSubmitting ? "Adding..." : "Add"}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

        </div>
    );
};

export default DatasetDetails;
