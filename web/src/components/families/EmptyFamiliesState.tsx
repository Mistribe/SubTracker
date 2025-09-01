import { CreateFamilyDialog } from "./CreateFamilyDialog";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Activity, Link2, PieChart, ShieldCheck, Share2, Users, Eye } from "lucide-react";

export const EmptyFamiliesState = () => {
  return (
    <div className="mx-auto max-w-3xl">
      <Card>
        <CardHeader className="text-center">
          <div className="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary">
            <Users className="h-5 w-5" />
          </div>
          <CardTitle className="text-2xl">Create your family</CardTitle>
          <CardDescription className="text-base">
            Families in SubTracker help you manage subscriptions together with the people who matter.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ul className="grid gap-4 sm:grid-cols-2">
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary"><Share2 className="h-5 w-5" /></div>
              <div>
                <p className="font-medium">Share subscriptions</p>
                <p className="text-muted-foreground">Give access to plans your family can use together.</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary"><Eye className="h-5 w-5" /></div>
              <div>
                <p className="font-medium">See who uses what</p>
                <p className="text-muted-foreground">Know exactly which member uses each subscription.</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary"><PieChart className="h-5 w-5" /></div>
              <div>
                <p className="font-medium">Track spending</p>
                <p className="text-muted-foreground">View your personal spend and the family total at a glance.</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary"><Link2 className="h-5 w-5" /></div>
              <div>
                <p className="font-medium">Invite with a link</p>
                <p className="text-muted-foreground">Send secure invite links or emails to add members fast.</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary"><ShieldCheck className="h-5 w-5" /></div>
              <div>
                <p className="font-medium">Roles & permissions</p>
                <p className="text-muted-foreground">Control who can manage members, subscriptions, and billing.</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary"><Activity className="h-5 w-5" /></div>
              <div>
                <p className="font-medium">Activity & renewals</p>
                <p className="text-muted-foreground">See recent changes and upcoming renewals to avoid surprises.</p>
              </div>
            </li>
          </ul>
        </CardContent>
        <CardFooter className="justify-center">
          <CreateFamilyDialog />
        </CardFooter>
      </Card>
      <p className="mt-4 text-center text-xs text-muted-foreground">
        You can update family settings at any time.
      </p>
    </div>
  );
};