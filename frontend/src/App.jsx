import React from "react";
import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import MainNavbar from "./components/MainNavbar";
import Footer from "./components/Footer";
import routes from "./routes";
import TokenManager from "./services/TokenManager";

function App() {
    return (
        <React.StrictMode>
            <TokenManager />
            <BrowserRouter>
                <div className="main-container">
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
                    <Footer />
                </div>
            </BrowserRouter>
        </React.StrictMode>
    );
}

export default App;
