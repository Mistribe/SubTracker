# Project Guidelines for SubTracker

## Project Overview

SubTracker is a React TypeScript application built with Vite, focusing on subscription management functionality.
The project uses modern web development practices with a component-based architecture.

## Tech Stack

- **Frontend**: React 19.1.0 with TypeScript 5.8.3
- **Build Tool**: Vite 7.0.4
- **Styling**: Tailwind CSS 4.1.11 with shadcn/ui components ⚠️ **PRIORITY**
- **State Management**: React Query (@tanstack/react-query)
- **Routing**: React Router DOM
- **Authentication**: Kinde Auth React
- **Form Management**: React Hook Form with Zod validation
- **UI Components**: Radix UI primitives
- **Package Manager**: npm

## Styling and UI Priority Guidelines

**🎯 PRIORITY: Always use Tailwind CSS and shadcn/ui components first**

### Primary Styling Approach

1. **Tailwind CSS First**: Use Tailwind utility classes for all styling needs
    - Responsive design with Tailwind breakpoints (`sm:`, `md:`, `lg:`, `xl:`, `2xl:`)
    - Spacing, colors, typography, and layout using Tailwind classes
    - Custom CSS should be minimal and only when Tailwind utilities are insufficient

2. **shadcn/ui Components Priority**: Always check if a shadcn/ui component exists before creating custom components
    - Use components from `/src/components/ui/` directory
    - Available components include: Button, Input, Dialog, Card, Select, Checkbox, etc.
    - Follow shadcn/ui patterns for component composition and variants

3. **Radix UI Integration**: shadcn/ui components are built on Radix UI primitives
    - Inherit accessibility features from Radix UI automatically
    - Use Radix UI directly only when shadcn/ui doesn't provide the needed component
    - Maintain consistent behavior patterns across the application

### Component Development Hierarchy

1. **First**: Check if shadcn/ui has the component you need
2. **Second**: Use Radix UI primitives if shadcn/ui doesn't have it
3. **Third**: Create custom components using Tailwind classes
4. **Last Resort**: Write custom CSS (avoid when possible)

### Tailwind Configuration

- Use the configured Tailwind theme and design tokens
- Leverage CSS variables for theme switching (dark/light mode)
- Follow the established spacing, color, and typography scale
- Use Tailwind's aspect-ratio, grid, and flexbox utilities for layouts

### Design System Consistency

- Maintain visual consistency by using established component variants
- Follow the existing color palette and spacing system
- Use consistent border radius, shadows, and transitions
- Ensure all interactive elements follow the same hover/focus states

### 🎨 Styling Priority Guidelines

**ALWAYS use Tailwind CSS and shadcn/ui components as the primary styling solution:**

- ✅ **First Priority**: Use existing shadcn/ui components from `/src/components/ui/`
- ✅ **Second Priority**: Style with Tailwind CSS utility classes
- ❌ **Avoid**: Custom CSS files, inline styles, or other CSS frameworks
- ❌ **Avoid**: Writing custom components when shadcn/ui equivalents exist


