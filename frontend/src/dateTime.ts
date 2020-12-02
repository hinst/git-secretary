export function getStartOfDay(dateTime: Date) {
    return new Date(dateTime.getFullYear(), dateTime.getMonth(), dateTime.getDate());
}