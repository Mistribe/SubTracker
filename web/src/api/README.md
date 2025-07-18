# SubTracker API Client

This directory contains the API client and related utilities for interacting with the SubTracker backend API.

## Overview

The API client is built using Axios and provides typed methods for all endpoints defined in the SubTracker API. It includes proper error handling and response typing based on the Swagger definition.

## Usage

### Basic Usage

```typescript
import { api } from '../api';

// Example: Get all subscriptions
const fetchSubscriptions = async () => {
  try {
    const subscriptions = await api.getSubscriptions();
    // Do something with the subscriptions
    console.log(subscriptions);
  } catch (error) {
    console.error('Failed to fetch subscriptions:', error);
  }
};
```

### Error Handling

The API client includes built-in error handling, but you can also use the `handleApiError` utility function for component-level error handling:

```typescript
import { api, handleApiError } from '../api';
import { useState } from 'react';

const SubscriptionList = () => {
  const [error, setError] = useState<string | null>(null);
  
  const fetchSubscriptions = async () => {
    try {
      const subscriptions = await api.getSubscriptions();
      // Do something with the subscriptions
    } catch (error) {
      setError(handleApiError(error, 'Failed to load subscriptions'));
    }
  };
  
  return (
    <div>
      {error && <div className="error">{error}</div>}
      {/* Rest of the component */}
    </div>
  );
};
```

### Utility Functions

The API client comes with several utility functions to help with common tasks:

- `formatDate(dateString)`: Formats a date string for display
- `formatCurrency(amount, currency)`: Formats a currency value for display
- `calculateMonthlyCost(price, months)`: Calculates the monthly cost of a subscription
- `isSubscriptionActive(endDate)`: Checks if a subscription is active
- `getTimeRemaining(endDate)`: Gets the time remaining until a subscription expires

Example:

```typescript
import { formatCurrency, formatDate, isSubscriptionActive } from '../api';

const SubscriptionCard = ({ subscription }) => {
  const latestPayment = subscription.payments[0];
  const isActive = isSubscriptionActive(latestPayment.end_date);
  
  return (
    <div className={`card ${isActive ? 'active' : 'expired'}`}>
      <h3>{subscription.name}</h3>
      <p>Price: {formatCurrency(latestPayment.price, latestPayment.currency)}</p>
      <p>Expires: {formatDate(latestPayment.end_date)}</p>
    </div>
  );
};
```

## Models

The API client uses TypeScript interfaces to provide type safety when working with the API. These models are defined in the `models` directory and are based on the Swagger definition.

### Available Models

- **Label**: `LabelModel`, `CreateLabelModel`, `UpdateLabelModel`
- **Family**: `FamilyModel`, `CreateFamilyModel`, `UpdateFamilyModel`
- **Family Member**: `FamilyMemberModel`, `CreateFamilyMemberModel`, `UpdateFamilyMemberModel`
- **Subscription**: `SubscriptionModel`, `CreateSubscriptionModel`, `UpdateSubscriptionModel`
- **Payment**: `PaymentModel`, `CreatePaymentModel`, `CreateSubscriptionPaymentModel`
- **Error**: `HttpError`

Example:

```typescript
import { CreateSubscriptionModel } from '../api';

// Create a new subscription
const newSubscription: CreateSubscriptionModel = {
  name: 'Netflix',
  family_id: 'family-123',
  family_members: ['member-1', 'member-2'],
  labels: ['streaming', 'entertainment'],
  payer_id: 'member-1',
  payed_by_joint_account: false,
  payments: [
    {
      price: 15.99,
      currency: 'USD',
      months: 1,
      start_date: '2025-07-18',
      end_date: '2025-08-18'
    }
  ]
};

await api.createSubscription(newSubscription);
```

## API Endpoints

The API client provides methods for all endpoints defined in the SubTracker API:

### Label Endpoints
- `getLabels(withDefault?)`: Get all labels
- `getDefaultLabels()`: Get default labels
- `getLabelById(id)`: Get label by ID
- `createLabel(label)`: Create a new label
- `updateLabel(id, label)`: Update a label
- `deleteLabel(id)`: Delete a label

### Family Endpoints
- `getFamilies()`: Get all families
- `getFamilyById(id)`: Get family by ID
- `createFamily(family)`: Create a new family
- `updateFamily(id, family)`: Update a family
- `addFamilyMember(familyId, member)`: Add a family member
- `updateFamilyMember(familyId, memberId, member)`: Update a family member
- `deleteFamilyMember(familyId, memberId)`: Delete a family member

### Subscription Endpoints
- `getSubscriptions()`: Get all subscriptions
- `getSubscriptionById(id)`: Get subscription by ID
- `createSubscription(subscription)`: Create a new subscription
- `updateSubscription(id, subscription)`: Update a subscription
- `deleteSubscription(id)`: Delete a subscription
- `addSubscriptionPayment(subscriptionId, payment)`: Add a payment to a subscription
- `deleteSubscriptionPayment(subscriptionId, paymentId)`: Delete a subscription payment