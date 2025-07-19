import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import type {ReactNode} from "react";

interface ProtectedRouteProps {
    children: ReactNode;
}

export const ProtectedRoute = ({children}: ProtectedRouteProps) => {
    const {isAuthenticated, login} = useKindeAuth();

    if (!isAuthenticated) {

        login();
        return null;
    }

    return <>{children}</>;
};