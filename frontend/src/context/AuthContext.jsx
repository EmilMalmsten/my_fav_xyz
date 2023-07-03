import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';

export const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [authUser, setAuthUser] = useState(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  const login = async (email, password) => {
    try {
        const response = await axios.post(`${import.meta.env.VITE_API_URL}/login`, {
          email,
          password,
        });
        console.log(response)
        const user = {
            email: response.data.user.email,
            userID: response.data.user.id,
            createdAt: response.data.user.created_at,
        }
        localStorage.setItem('accessToken', response.data.token);
        localStorage.setItem('refreshToken', response.data.refresh_token);
        localStorage.setItem('authUser', user);
        setIsLoggedIn(true);
        setAuthUser(user)  
    } catch (e) {
        console.log(e);
    }
  };

  const logout = async () => {
    try {
        const refreshToken = localStorage.getItem('refreshToken');
        await axios.post(`${import.meta.env.VITE_API_URL}/revoke`, null, {
            headers: {
              Authorization: `Bearer ${refreshToken}`,
            },
        });
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('user');
        setIsLoggedIn(false);
        setAuthUser(null);
    } catch (e) {
        console.log(e);
    }
  };

  useEffect(() => {
    const storedAuthUser = localStorage.getItem('authUser');
    const storedIsLoggedIn = localStorage.getItem('isLoggedIn');

    if (storedAuthUser && storedIsLoggedIn) {
      setAuthUser(JSON.parse(storedAuthUser));
      setIsLoggedIn(storedIsLoggedIn === 'true');
    }
  }, []);

  const value = {
    authUser,
    login,
    isLoggedIn,
    logout
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
