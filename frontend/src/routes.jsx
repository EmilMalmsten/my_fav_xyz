import PropTypes from "prop-types";
import { Navigate } from "react-router-dom";
import Home from "./pages/Home";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Toplist from "./pages/ViewToplist";
import CreateToplist from "./pages/CreateToplist";
import ToplistsByCategory from "./components/ToplistsByCategory";
import { useAuth } from "./context/AuthContext";

const ProtectedRoute = ({ element }) => {
    const { isLoggedIn } = useAuth();

    if (!isLoggedIn) {
        return <Navigate to="/login" />;
    }

    return element;
};

ProtectedRoute.propTypes = {
    path: PropTypes.string.isRequired,
    element: PropTypes.element.isRequired,
};

const routes = [
    {
        path: "/",
        element: <Home />,
    },
    {
        path: "/register",
        element: <Register />,
    },
    {
        path: "/login",
        element: <Login />,
    },
    {
        path: "/toplists/:id",
        element: <ProtectedRoute element={<Toplist />} />,
    },
    {
        path: "/toplists/create",
        element: <ProtectedRoute element={<CreateToplist />} />,
    },
    {
        path: "/toplists/recent",
        element: (
            <ToplistsByCategory
                title="Most recent toplists"
                endpoint="/toplists/recent"
            />
        ),
    },
    {
        path: "/toplists/popular",
        element: (
            <ToplistsByCategory
                title="Most popular toplists"
                endpoint="/toplists/popular"
            />
        ),
    },
    {
        path: "*",
        element: <Navigate to="/" />,
    },
];

export default routes;
