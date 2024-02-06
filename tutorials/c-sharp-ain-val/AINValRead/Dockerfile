FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /app

# Copy everything
COPY . ./

# Publish the application
RUN dotnet publish "ReadAIN.csproj" -c Release -o /app

# Set the working directory for the final image
WORKDIR /app

# Set the entry point
ENTRYPOINT ["dotnet", "ReadAIN.dll"]
