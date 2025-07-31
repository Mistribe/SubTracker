import {useState} from "react";
import {Button} from "@/components/ui/button";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Edit2, Plus, Trash2} from "lucide-react";
import Plan from "@/models/plan";
import Price from "@/models/price";
import {PlanForm, type PlanFormValues} from "./PlanForm";
import {PriceForm, type PriceFormValues} from "./PriceForm";

interface PlansListProps {
    plans: Plan[];
    currencies: string[];
    isLoadingCurrencies: boolean;
    onAddPlan: (data: PlanFormValues) => Promise<void>;
    onUpdatePlan: (planId: string, data: PlanFormValues) => Promise<void>;
    onDeletePlan: (planId: string) => Promise<void>;
    onAddPrice: (planId: string, data: PriceFormValues) => Promise<void>;
    onUpdatePrice: (planId: string, priceId: string, data: PriceFormValues) => Promise<void>;
    onDeletePrice: (planId: string, priceId: string) => Promise<void>;
    isSubmitting: boolean;
}

export function PlansList({
                              plans,
                              currencies,
                              isLoadingCurrencies,
                              onAddPlan,
                              onUpdatePlan,
                              onDeletePlan,
                              onAddPrice,
                              onUpdatePrice,
                              onDeletePrice,
                              isSubmitting
                          }: PlansListProps) {
    // State for managing UI
    const [showAddPlanForm, setShowAddPlanForm] = useState(false);
    const [editingPlanId, setEditingPlanId] = useState<string | null>(null);
    const [showAddPriceForm, setShowAddPriceForm] = useState<string | null>(null); // planId or null
    const [editingPriceInfo, setEditingPriceInfo] = useState<{ planId: string, priceId: string } | null>(null);

    // Handle plan form submission
    const handlePlanSubmit = async (data: PlanFormValues) => {
        if (editingPlanId) {
            await onUpdatePlan(editingPlanId, data);
        } else {
            await onAddPlan(data);
        }
        setShowAddPlanForm(false);
        setEditingPlanId(null);
    };

    // Handle price form submission
    const handlePriceSubmit = async (data: PriceFormValues) => {
        if (editingPriceInfo) {
            await onUpdatePrice(editingPriceInfo.planId, editingPriceInfo.priceId, data);
            setEditingPriceInfo(null);
        } else if (showAddPriceForm) {
            await onAddPrice(showAddPriceForm, data);
        }
        setShowAddPriceForm(null);
    };

    // Start editing a plan
    const startEditingPlan = (plan: Plan) => {
        setEditingPlanId(plan.id);
        setShowAddPlanForm(true);
    };

    // Start editing a price
    const startEditingPrice = (planId: string, price: Price) => {
        setEditingPriceInfo({planId, priceId: price.id});
        setShowAddPriceForm(planId);
    };

    // Find the plan being edited
    const editingPlan = editingPlanId ? plans.find(p => p.id === editingPlanId) : undefined;

    // Find the price being edited
    const editingPrice = editingPriceInfo
        ? plans.find(p => p.id === editingPriceInfo.planId)?.prices.find(p => p.id === editingPriceInfo.priceId)
        : undefined;

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h3 className="text-lg font-medium">Plans</h3>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => {
                        setEditingPlanId(null);
                        setShowAddPlanForm(true);
                    }}
                >
                    <Plus className="h-4 w-4 mr-2"/>
                    Add Plan
                </Button>
            </div>

            {/* Add/Edit Plan Form */}
            {showAddPlanForm && (
                <PlanForm
                    plan={editingPlan}
                    onSubmit={handlePlanSubmit}
                    onCancel={() => {
                        setShowAddPlanForm(false);
                        setEditingPlanId(null);
                    }}
                    isSubmitting={isSubmitting}
                />
            )}

            {/* Plans list */}
            {plans.length === 0 && !showAddPlanForm ? (
                <div className="text-center py-4 text-gray-500">
                    No plans added yet. Click "Add Plan" to create one.
                </div>
            ) : (
                <div className="space-y-4 mt-4">
                    {plans.map((plan) => (
                        <Card key={plan.id}>
                            <CardHeader className="pb-2">
                                <div className="flex justify-between items-start">
                                    <CardTitle>{plan.name}</CardTitle>
                                    <div className="flex space-x-2">
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => startEditingPlan(plan)}
                                        >
                                            <Edit2 className="h-4 w-4"/>
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => onDeletePlan(plan.id)}
                                        >
                                            <Trash2 className="h-4 w-4 text-red-500"/>
                                        </Button>
                                    </div>
                                </div>
                                {plan.description && (
                                    <p className="text-sm text-gray-500">{plan.description}</p>
                                )}
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-2">
                                    <div className="flex justify-between items-center">
                                        <h4 className="text-sm font-medium">Prices</h4>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => {
                                                setEditingPriceInfo(null);
                                                setShowAddPriceForm(plan.id);
                                            }}
                                        >
                                            <Plus className="h-4 w-4 mr-1"/>
                                            Add Price
                                        </Button>
                                    </div>

                                    {/* Add/Edit Price Form */}
                                    {showAddPriceForm === plan.id && (
                                        <PriceForm
                                            price={editingPrice}
                                            existingPrices={plan.prices}
                                            currencies={currencies}
                                            isLoadingCurrencies={isLoadingCurrencies}
                                            onSubmit={handlePriceSubmit}
                                            onCancel={() => {
                                                setShowAddPriceForm(null);
                                                setEditingPriceInfo(null);
                                            }}
                                            isSubmitting={isSubmitting}
                                        />
                                    )}

                                    {plan.prices.length === 0 && showAddPriceForm !== plan.id ? (
                                        <div className="text-center py-2 text-gray-500 text-sm">
                                            No prices added yet.
                                        </div>
                                    ) : (
                                        <div className="space-y-2">
                                            {plan.prices.map((price) => (
                                                <div key={price.id}
                                                     className="flex justify-between items-center p-2 border rounded-md">
                                                    <div>
                                                        <div className="font-medium">
                                                            {price.amount} {price.currency}
                                                        </div>
                                                        <div className="text-xs text-gray-500">
                                                            From: {price.startDate.toLocaleDateString()}
                                                            {price.endDate && ` To: ${price.endDate.toLocaleDateString()}`}
                                                        </div>
                                                    </div>
                                                    <div className="flex space-x-1">
                                                        <Button
                                                            variant="ghost"
                                                            size="sm"
                                                            onClick={() => startEditingPrice(plan.id, price)}
                                                        >
                                                            <Edit2 className="h-3 w-3"/>
                                                        </Button>
                                                        <Button
                                                            variant="ghost"
                                                            size="sm"
                                                            onClick={() => onDeletePrice(plan.id, price.id)}
                                                        >
                                                            <Trash2 className="h-3 w-3 text-red-500"/>
                                                        </Button>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}