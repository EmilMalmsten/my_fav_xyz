import { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";
import PropTypes from "prop-types";

export const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [authUser, setAuthUser] = useState(null);
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    const login = async (email, password) => {
        const response = await axios.post(
            `${import.meta.env.VITE_API_URL}/login`,
            {
                email,
                password,
            }
        );
        const user = {
            email: response.data.user.email,
            userID: response.data.user.id,
            createdAt: response.data.user.created_at,
        };
        localStorage.setItem("accessToken", response.data.token);
        localStorage.setItem("refreshToken", response.data.refresh_token);
        localStorage.setItem("authUser", JSON.stringify(user));
        setIsLoggedIn(true);
        setAuthUser(user);
    };

    const logout = async () => {
        try {
            const refreshToken = localStorage.getItem("refreshToken");
            await axios.post(`${import.meta.env.VITE_API_URL}/revoke`, null, {
                headers: {
                    Authorization: `Bearer ${refreshToken}`,
                },
            });
            localStorage.removeItem("accessToken");
            localStorage.removeItem("refreshToken");
            localStorage.removeItem("authUser");
            setIsLoggedIn(false);
            setAuthUser(null);
        } catch (e) {
            console.log(e);
        }
    };

    useEffect(() => {
        const storedAuthUser = localStorage.getItem("authUser");
        const storedAccessToken = localStorage.getItem("accessToken");
        const storedRefreshToken = localStorage.getItem("refreshToken");

        if (storedAuthUser && storedAccessToken && storedRefreshToken) {
            const authUser = JSON.parse(storedAuthUser);
            setAuthUser(authUser);
            setIsLoggedIn(true);
        } else {
            setAuthUser(null);
            setIsLoggedIn(false);
        }
    }, []);

    const value = {
        authUser,
        login,
        isLoggedIn,
        logout,
    };

    AuthProvider.propTypes = {
        children: PropTypes.node.isRequired,
    };

    return (
        <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    );
};
