// Common interfaces for API responses and requests

// HTTP Error model
export interface HttpError {
  message: string;
}

// Label models
export interface LabelModel {
  id: string;
  name: string;
  color: string;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateLabelModel {
  id?: string;
  name: string;
  color: string;
  is_default?: boolean;
  created_at?: string;
}

export interface UpdateLabelModel {
  name?: string;
  color?: string;
  updated_at?: string;
}

// Family models
export interface FamilyMemberModel {
  id: string;
  family_id: string;
  name: string;
  is_kid: boolean;
  created_at: string;
  updated_at: string;
}

export interface FamilyModel {
  id: string;
  name: string;
  have_joint_account: boolean;
  is_owner: boolean;
  members: FamilyMemberModel[];
  created_at: string;
  updated_at: string;
}

export interface CreateFamilyModel {
  id?: string;
  name: string;
  have_joint_account: boolean;
  created_at?: string;
}

export interface UpdateFamilyModel {
  name?: string;
  have_joint_account?: boolean;
  updated_at?: string;
}

export interface CreateFamilyMemberModel {
  id?: string;
  name: string;
  is_kid: boolean;
  email?: string;
  created_at?: string;
}

export interface UpdateFamilyMemberModel {
  name?: string;
  id_kid?: boolean; // Note: This matches the Swagger definition, but might be a typo (should be is_kid)
  email?: string;
  updated_at?: string;
}

// Payment models
export interface PaymentModel {
  id: string;
  price: number;
  currency: string;
  months: number;
  start_date: string;
  end_date: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePaymentModel {
  id?: string;
  price: number;
  currency: string;
  months: number;
  start_date: string;
  end_date: string;
  created_at?: string;
}

export type CreateSubscriptionPaymentModel = CreatePaymentModel;

// Subscription models
export interface SubscriptionModel {
  id: string;
  name: string;
  family_members: string[];
  labels: string[];
  payer: string;
  payments: PaymentModel[];
  created_at: string;
  updated_at: string;
}

export interface CreateSubscriptionModel {
  id?: string;
  name: string;
  family_id: string;
  family_members: string[];
  labels: string[];
  payer_id: string;
  payed_by_joint_account: boolean;
  payments: CreatePaymentModel[];
  created_at?: string;
}

export interface UpdateSubscriptionModel {
  name?: string;
  family_members?: string[];
  labels?: string[];
  payer_id?: string;
  payed_by_joint_account?: boolean;
  updated_at?: string;
}