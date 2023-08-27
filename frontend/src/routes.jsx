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
import UpdateUserEmail from "./pages/UpdateUserEmail";
import UserDashboard from "./pages/UserDashboard";
import UpdateUserPassword from "./pages/UpdateUserPassword";
import DeleteUser from "./pages/DeleteUser";
import ToplistSearchResults from "./pages/ToplistSearchResults";

const ProtectedRoute = ({ element }) => {
    const { getLoginStatus } = useAuth();
    const loggedIn = getLoginStatus();

    if (loggedIn) {
        return element;
    } else {
        return <Navigate to="/login" />;
    }
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
        path: "/toplists/search",
        element: <ToplistSearchResults />,
    },
    {
        path: "/toplists/create",
        element: <ProtectedRoute element={<CreateToplist />} />,
    },
    {
        path: "/users/:id/email",
        element: <ProtectedRoute element={<UpdateUserEmail />} />,
    },
    {
        path: "/users/:id/password",
        element: <ProtectedRoute element={<UpdateUserPassword />} />,
    },
    {
        path: "/users/:id/delete",
        element: <ProtectedRoute element={<DeleteUser />} />,
    },
    {
        path: "/users/:id",
        element: <ProtectedRoute element={<UserDashboard />} />,
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
