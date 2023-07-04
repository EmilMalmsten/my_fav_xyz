import React from "react";
import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import MainNavbar from "./components/MainNavbar";
import routes from "./routes";
import TokenManager from "./services/TokenManager";

function App() {
    return (
        <React.StrictMode>
            <TokenManager />
            <BrowserRouter>
                <MainNavbar />
                <Routes>
                    {routes.map((route, index) => (
                        <Route
                            key={index}
                            path={route.path}
                            element={route.element}
                        />
                    ))}
                </Routes>
            </BrowserRouter>
        </React.StrictMode>
    );
}

export default App;
