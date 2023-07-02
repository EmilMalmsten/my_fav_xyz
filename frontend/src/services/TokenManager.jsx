import { useEffect } from 'react';
import jwtDecode from 'jwt-decode';
import axios from 'axios';

function TokenManager() {
    const checkTokenExpiry = () => {
        const accessToken = localStorage.getItem('accessToken');
        if (!accessToken || accessToken === 'undefined' || accessToken === 'null') {
            console.log('no access token')
            return;
        }
      
        const decodedToken = decodeAccessToken(accessToken);
        const currentTime = Math.floor(Date.now() / 1000); // ms to s
        if (decodedToken.exp < currentTime) {
            console.log('token expired, refreshing token')
            const refreshToken = localStorage.getItem('refreshToken');
            refreshAccessToken(refreshToken);
            return;
        }
      
        console.log('not expired!')
      };

  const decodeAccessToken = (accessToken) => {
    try {
      const decodedToken = jwtDecode(accessToken);
      return decodedToken;
    } catch (error) {
      console.error('Error decoding access token:', error);
      return null;
    }
  };

  const refreshAccessToken = async (refreshToken) => {
    try {
        const response = await axios.post(`${import.meta.env.VITE_API_URL}/refresh`, null, {
          headers: {
            Authorization: `Bearer ${refreshToken}`,
          },
        });
  
        const { token: newAccessToken } = response.data;
        localStorage.setItem('accessToken', newAccessToken);
        console.log('token refreshed');
        return newAccessToken;
      } catch (error) {
        console.error('Error refreshing token:', error);
        throw error;
      }
  };

  useEffect(() => {
    checkTokenExpiry();
    const interval = setInterval(checkTokenExpiry, 20000); // 1 minute interval

    return () => clearInterval(interval);
  }, []);

  return null;
}

export default TokenManager;
