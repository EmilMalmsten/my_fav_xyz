import { useEffect } from "react";
import jwtDecode from "jwt-decode";
import axios from "axios";
import { useAuth } from "../context/AuthContext";

function TokenManager() {
    const { logout } = useAuth();

    const checkAccessTokenExpiry = () => {
        const accessToken = localStorage.getItem("accessToken");
        if (
            !accessToken ||
            accessToken === "undefined" ||
            accessToken === "null"
        ) {
            return;
        }

        const decodedToken = decodeJWT(accessToken);
        const currentTime = Math.floor(Date.now() / 1000); // ms to s
        if (decodedToken.exp < currentTime) {
            const refreshToken = localStorage.getItem("refreshToken");
            refreshAccessToken(refreshToken);
            return;
        }
    };

    const checkRefreshTokenExpiry = () => {
        const refreshToken = localStorage.getItem("refreshToken");
        if (
            !refreshToken ||
            refreshToken === "undefined" ||
            refreshToken === "null"
        ) {
            return;
        }

        const decodedToken = decodeJWT(refreshToken);
        const currentTime = Math.floor(Date.now() / 1000); // ms to s
        if (decodedToken.exp < currentTime) {
            logout();
            return;
        }
    };

    const decodeJWT = (token) => {
        try {
            const decodedToken = jwtDecode(token);
            return decodedToken;
        } catch (error) {
            console.error("Error decoding access token:", error);
            return null;
        }
    };

    const refreshAccessToken = async (refreshToken) => {
        try {
            const response = await axios.post(
                `${import.meta.env.VITE_API_URL}/refresh`,
                null,
                {
                    headers: {
                        Authorization: `Bearer ${refreshToken}`,
                    },
                }
            );

            const { token: newAccessToken } = response.data;
            localStorage.setItem("accessToken", newAccessToken);
            return newAccessToken;
        } catch (error) {
            console.error("Error refreshing token:", error);
            throw error;
        }
    };

    useEffect(() => {
        checkAccessTokenExpiry();
        checkRefreshTokenExpiry();
        const interval = setInterval(checkAccessTokenExpiry, 20000); // 1 minute interval

        return () => clearInterval(interval);
    }, []);

    return null;
}

export default TokenManager;
