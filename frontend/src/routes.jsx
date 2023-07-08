import PropTypes from "prop-types";
import { Navigate } from "react-router-dom";
import Home from "./pages/Home";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Toplist from "./pages/ViewToplist";
import CreateToplist from "./pages/CreateToplist";
import ToplistsByCategory from "./components/ToplistsByCategory";
import { useAuth } from "./context/AuthContext";
import EditToplist from "./pages/EditToplist";
import EditToplistItems from "./pages/EditToplistItems";

const ProtectedRoute = ({ element }) => {
    const { isLoggedIn } = useAuth();

    if (!isLoggedIn) {
        return <Navigate to="/login" />;
    }

    return element;
};

ProtectedRoute.propTypes = {
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
        element: <Toplist />,
    },
    {
        path: "/toplists/create",
        element: <ProtectedRoute element={<CreateToplist />} />,
    },
    {
        path: "/toplists/:id/edit",
        element: <ProtectedRoute element={<EditToplist />} />,
        children: [
            {
                path: "",
                element: <EditToplist />,
            },
        ],
    },
    {
        path: "/toplists/:id/items",
        element: <ProtectedRoute element={<EditToplistItems />} />,
        children: [
            {
                path: "",
                element: <EditToplistItems />,
            },
        ],
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
