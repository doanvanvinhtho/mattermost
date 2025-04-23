# MM-63240 Implementation Plan: TeamSettings.BrowseArchivedPublicChannels

## Overview

This issue involves replacing the `TeamSettings.ExperimentalViewArchivedChannels` setting with a new, more focused setting called `TeamSettings.BrowseArchivedPublicChannels`. This change provides more granular control over the visibility of archived channels.

## Current Behavior

- `TeamSettings.ExperimentalViewArchivedChannels` controls whether users can access and see archived channels.
- When enabled (true), users can access archived channels and see them listed in places like the channel switcher.
- When disabled (false), users cannot access or see archived channels at all.
- This setting has defaulted to true since October 2022 (v5.28).

## Desired Behavior

- A new setting `TeamSettings.BrowseArchivedPublicChannels` will default to true.
- This setting will specifically control the visibility of archived public channels when browsing available channels:
  - When enabled (true): Archived public channels will be visible when browsing available channels
  - When disabled (false): Archived public channels will be hidden when browsing available channels, unless the user is already a member
- Users will always be able to access archived channels where they are members.
- Users can leave archived channels if they no longer wish to see them.

## Implementation Plan

### 1. Add the New Setting in Config

- Add `BrowseArchivedPublicChannels` to the `TeamSettings` struct in `/public/model/config.go`
- Set a default value of `true` in the `SetDefaults()` method
- Add appropriate JSON tags and documentation

### 2. Update API Endpoints for Channel Browsing

- Modify the `getPublicChannelsForTeam` endpoint in `/channels/api4/channel.go`
- Pass the new setting to the App layer

### 3. Update App Layer

- Modify the `GetPublicChannelsForTeam` method in `/channels/app/channel.go`
- Pass the new setting to the Store layer

### 4. Update Store Layer

- Modify the `GetPublicChannelsForTeam` method in `/channels/store/sqlstore/channel_store.go`
- Update the query to conditionally include/exclude archived channels based on the new setting

### 5. Update Authorization Logic

- Ensure that users can still access archived channels where they are members
- Update any permission checks to use the new setting

### 6. Update Client Config

- Expose the new setting in `/config/client.go` for the frontend to use

### 7. Add Migration Logic

- Add logic to migrate existing `ExperimentalViewArchivedChannels` settings to the new `BrowseArchivedPublicChannels` setting
- Handle the case where the old setting is false (new setting should also be false)

### 8. Add Telemetry

- Update telemetry reporting to include the new setting

### 9. Add Tests

- Add unit tests for the new setting
- Update existing tests that use `ExperimentalViewArchivedChannels`
- Add integration tests for the new behavior

### 10. Maintain Backward Compatibility (Optional)

- Keep the old setting for a transition period
- Add deprecation warnings for the old setting

## Notes

- The old setting controlled visibility in many places, while the new setting specifically focuses on browsing public channels
- This change aims to provide more granular control while maintaining a consistent user experience
- We'll need to ensure that users can still access archived channels where they are members, regardless of the new setting's value

## Technical Concerns

- We need to ensure we don't break existing functionality
- We should consider the impact on API clients that might be using the current behavior
- We should document the change clearly for administrators